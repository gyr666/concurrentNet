#!/usr/bin/env ruby

# TODO (temporary here, we'll move this into the Github issues once
#       redis-trib initial implementation is completed).
#
# - Make sure that if the rehashing fails in the middle redis-trib will try
#   to recover.
# - When redis-trib performs a cluster check, if it detects a slot move in
#   progress it should prompt the user to continue the move from where it
#   stopped.
# - Gracefully handle Ctrl+C in move_slot to prompt the user if really stop
#   while rehashing, and performing the best cleanup possible if the user
#   forces the quit.
# - When doing "fix" set a global Fix to true, and prompt the user to
#   fix the problem if automatically fixable every time there is something
#   to fix. For instance:
#   1) If there is a node that pretend to receive a slot, or to migrate a
#      slot, but has no entries in that slot, fix it.
#   2) If there is a node having keys in slots that are not owned by it
#      fix this condition moving the entries in the same node.
#   3) Perform more possibly slow tests about the state of the cluster.
#   4) When aborted slot migration is detected, fix it.

require 'rubygems'
require 'redis'

ClusterHashSlots = 16384

def xputs(s)
    case s[0..2]
    when ">>>"
        color="29;1"
    when "[ER"
        color="31;1"
    when "[OK"
        color="32"
    when "[FA","***"
        color="33"
    else
        color=nil
    end

    color = nil if ENV['TERM'] != "xterm"
    print "\033[#{color}m" if color
    print s
    print "\033[0m" if color
    print "\n"
end

class ClusterNode
    def initialize(addr)
        s = addr.split(":")
        if s.length != 2
            puts "Invalid node name #{addr}"
            exit 1
        end
        @r = nil
        @info = {}
        @info[:host] = s[0]
        @info[:port] = s[1]
        @info[:slots] = {}
        @info[:migrating] = {}
        @info[:importing] = {}
        @info[:replicate] = false
        @dirty = false # True if we need to flush slots info into node.
        @friends = []
    end

    def friends
        @friends
    end

    def slots 
        @info[:slots]
    end

    def has_flag?(flag)
        @info[:flags].index(flag)
    end

    def to_s
        "#{@info[:host]}:#{@info[:port]}"
    end

    def connect(o={})
        return if @r
        print "Connecting to node #{self}: "
        STDOUT.flush
        begin
            @r = Redis.new(:host => @info[:host], :port => @info[:port])
            @r.ping
        rescue
            xputs "[ERR] Sorry, can't connect to node #{self}"
            exit 1 if o[:abort]
            @r = nil
        end
        xputs "OK"
    end

    def assert_cluster
        info = @r.info
        if !info["cluster_enabled"] || info["cluster_enabled"].to_i == 0
            xputs "[ERR] Node #{self} is not configured as a cluster node."
            exit 1
        end
    end

    def assert_empty
        if !(@r.cluster("info").split("\r\n").index("cluster_known_nodes:1")) ||
            (@r.info['db0'])
            xputs "[ERR] Node #{self} is not empty. Either the node already knows other nodes (check with CLUSTER NODES) or contains some key in database 0."
            exit 1
        end
    end

    def load_info(o={})
        self.connect
        nodes = @r.cluster("nodes").split("\n")
        nodes.each{|n|
            # name addr flags role ping_sent ping_recv link_status slots
            split = n.split
            name,addr,flags,role,ping_sent,ping_recv,config_epoch,link_status = split[0..6]
            slots = split[8..-1]
            info = {
                :name => name,
                :addr => addr,
                :flags => flags.split(","),
                :role => role,
                :ping_sent => ping_sent.to_i,
                :ping_recv => ping_recv.to_i,
                :link_status => link_status
            }
            if info[:flags].index("myself")
                @info = @info.merge(info)
                @info[:slots] = {}
                slots.each{|s|
                    if s[0..0] == '['
                        if s.index("->-") # Migrating
                            slot,dst = s[1..-1].split("->-")
                            @info[:migrating][slot] = dst
                        elsif s.index("-<-") # Importing
                            slot,src = s[1..-1].split("-<-")
                            @info[:importing][slot] = src
                        end
                    elsif s.index("-")
                        start,stop = s.split("-")
                        self.add_slots((start.to_i)..(stop.to_i))
                    else
                        self.add_slots((s.to_i)..(s.to_i))
                    end
                } if slots
                @dirty = false
                @r.cluster("info").split("\n").each{|e|    
                    k,v=e.split(":")
                    k = k.to_sym
                    v.chop!
                    if k != :cluster_state
                        @info[k] = v.to_i
                    else
                        @info[k] = v
                    end
                }
            elsif o[:getfriends]
                @friends << info
            end
        }
    end

    def add_slots(slots)
        slots.each{|s|
            @info[:slots][s] = :new
        }
        @dirty = true
    end

    def set_as_replica(node_id)
        @info[:replicate] = node_id
        @dirty = true
    end

    def flush_node_config
        return if !@dirty
        if @info[:replicate]
            begin
                @r.cluster("replicate",@info[:replicate])
            rescue
                # If the cluster did not already joined it is possible that
                # the slave does not know the master node yet. So on errors
                # we return ASAP leaving the dirty flag set, to flush the
                # config later.
                return
            end
        else
            new = []
            @info[:slots].each{|s,val|
                if val == :new
                    new << s
                    @info[:slots][s] = true
                end
            }
            @r.cluster("addslots",*new)
        end
        @dirty = false
    end

    def info_string
        # We want to display the hash slots assigned to this node
        # as ranges, like in: "1-5,8-9,20-25,30"
        #
        # Note: this could be easily written without side effects,
        # we use 'slots' just to split the computation into steps.
        
        # First step: we want an increasing array of integers
        # for instance: [1,2,3,4,5,8,9,20,21,22,23,24,25,30]
        slots = @info[:slots].keys.sort

        # As we want to aggregate adjacent slots we convert all the
        # slot integers into ranges (with just one element)
        # So we have something like [1..1,2..2, ... and so forth.
        slots.map!{|x| x..x}

        # Finally we group ranges with adjacent elements.
        slots = slots.reduce([]) {|a,b|
            if !a.empty? && b.first == (a[-1].last)+1
                a[0..-2] + [(a[-1].first)..(b.last)]
            else
                a + [b]
            end
        }

        # Now our task is easy, we just convert ranges with just one
        # element into a number, and a real range into a start-end format.
        # Finally we join the array using the comma as separator.
        slots = slots.map{|x|
            x.count == 1 ? x.first.to_s : "#{x.first}-#{x.last}"
        }.join(",")

        role = self.has_flag?("master") ? "M" : "S"

        if self.info[:replicate] and @dirty
            "S: #{self.info[:name]} #{self.to_s}"
        else
            "#{role}: #{self.info[:name]} #{self.to_s}\n"+
            "   slots:#{slots} (#{self.slots.length} slots) "+
            "#{(self.info[:flags]-["myself"]).join(",")}"
        end
    end

    # Return a single string representing nodes and associated slots.
    # TODO: remove slaves from config when slaves will be handled
    # by Redis Cluster.
    def get_config_signature
        config = []
        @r.cluster("nodes").each_line{|l|
            s = l.split
            slots = s[8..-1].select {|x| x[0..0] != "["}
            next if slots.length == 0
            config << s[0]+":"+(slots.sort.join(","))
        }
        config.sort.join("|")
    end

    def info
        @info
    end
    
    def is_dirty?
        @dirty
    end

    def r
        @r
    end
end

class RedisTrib
    def initialize
        @nodes = []
        @fix = false
        @errors = []
    end

    def check_arity(req_args, num_args)
        if ((req_args > 0 and num_args != req_args) ||
           (req_args < 0 and num_args < req_args.abs))
           xputs "[ERR] Wrong number of arguments for specified sub command"
           exit 1
        end
    end

    def add_node(node)
        @nodes << node
    end

    def cluster_error(msg)
        @errors << msg
        xputs msg
    end

    def get_node_by_name(name)
        @nodes.each{|n|
            return n if n.info[:name] == name.downcase
        }
        return nil
    end

    def check_cluster
        xputs ">>> Performing Cluster Check (using node #{@nodes[0]})"
        show_nodes
        check_config_consistency
        check_open_slots
        check_slots_coverage
    end

    # Merge slots of every known node. If the resulting slots are equal
    # to ClusterHashSlots, then all slots are served.
    def covered_slots
        slots = {}
        @nodes.each{|n|
            slots = slots.merge(n.slots)
        }
        slots
    end

    def check_slots_coverage
        xputs ">>> Check slots coverage..."
        slots = covered_slots
        if slots.length == ClusterHashSlots
            xputs "[OK] All #{ClusterHashSlots} slots covered."
        else
            cluster_error \
                "[ERR] Not all #{ClusterHashSlots} slots are covered by nodes."
            fix_slots_coverage if @fix
        end
    end

    def check_open_slots
        xputs ">>> Check for open slots..."
        open_slots = []
        @nodes.each{|n|
            if n.info[:migrating].size > 0
                cluster_error \
                    "[WARNING] Node #{n} has slots in migrating state."
                open_slots += n.info[:migrating].keys
            elsif n.info[:importing].size > 0
                cluster_error \
                    "[WARNING] Node #{n} has slots in importing state."
                open_slots += n.info[:importing].keys
            end
        }
        open_slots.uniq!
        if open_slots.length > 0
            xputs "[WARNING] The following slots are open: #{open_slots.join(",")}"
        end
        if @fix
            open_slots.each{|slot| fix_open_slot slot}
        end
    end

    def nodes_with_keys_in_slot(slot)
        nodes = []
        @nodes.each{|n|
            nodes << n if n.r.cluster("getkeysinslot",slot,1).length > 0
        }
        nodes
    end

    def fix_slots_coverage
        not_covered = (0...ClusterHashSlots).to_a - covered_slots.keys
        xputs ">>> Fixing slots coverage..."
        xputs "List of not covered slots: " + not_covered.join(",")

        # For every slot, take action depending on the actual condition:
        # 1) No node has keys for this slot.
        # 2) A single node has keys for this slot.
        # 3) Multiple nodes have keys for this slot.
        slots = {}
        not_covered.each{|slot|
            nodes = nodes_with_keys_in_slot(slot)
            slots[slot] = nodes
            xputs "Slot #{slot} has keys in #{nodes.length} nodes: #{nodes.join}"
        }

        none = slots.select {|k,v| v.length == 0}
        single = slots.select {|k,v| v.length == 1}
        multi = slots.select {|k,v| v.length > 1}

        # Handle case "1": keys in no node.
        if none.length > 0
            xputs "The folowing uncovered slots have no keys across the cluster:"
            xputs none.keys.join(",")
            yes_or_die "Fix these slots by covering with a random node?"
            none.each{|slot,nodes|
                node = @nodes.sample
                xputs ">>> Covering slot #{slot} with #{node}"
                node.r.cluster("addslots",slot)
            }
        end

        # Handle case "2": keys only in one node.
        if single.length > 0
            xputs "The folowing uncovered slots have keys in just one node:"
            puts single.keys.join(",")
            yes_or_die "Fix these slots by covering with those nodes?"
            single.each{|slot,nodes|
                xputs ">>> Covering slot #{slot} with #{nodes[0]}"
                nodes[0].r.cluster("addslots",slot)
            }
        end

        # Handle case "3": keys in multiple nodes.
        if multi.length > 0
            xputs "The folowing uncovered slots have keys in multiple nodes:"
            xputs multi.keys.join(",")
            yes_or_die "Fix these slots by moving keys into a single node?"
            multi.each{|slot,nodes|
                xputs ">>> Covering slot #{slot} moving keys to #{nodes[0]}"
                # TODO
                # 1) Set all nodes as "MIGRATING" for this slot, so that we
                # can access keys in the hash slot using ASKING.
                # 2) Move everything to node[0]
                # 3) Clear MIGRATING from nodes, and ADDSLOTS the slot to
                # node[0].
                raise "TODO: Work in progress"
            }
        end
    end

    # Slot 'slot' was found to be in importing or migrating state in one or
    # more nodes. This function fixes this condition by migrating keys where
    # it seems more sensible.
    def fix_open_slot(slot)
        migrating = []
        importing = []
        @nodes.each{|n|
            next if n.has_flag? "slave"
            if n.info[:migrating][slot]
                migrating << n
            elsif n.info[:importing][slot]
                importing << n
            elsif n.r.cluster("countkeysinslot",slot) > 0
                xputs "*** Found keys about slot #{slot} in node #{n}!"
            end
        }
        puts ">>> Fixing open slot #{slot}"
        puts "Set as migrating in: #{migrating.join(",")}"
        puts "Set as importing in: #{importing.join(",")}"

        # Case 1: The slot is in migrating state in one slot, and in
        #         importing state in 1 slot. That's trivial to address.
        if migrating.length == 1 && importing.length == 1
            move_slot(migrating[0],importing[0],slot,:verbose=>true)
        else
            xputs "[ERR] Sorry, Redis-trib can't fix this slot yet (work in progress)"
        end
    end

    # Check if all the nodes agree about the cluster configuration
    def check_config_consistency
        if !is_config_consistent?
            cluster_error "[ERR] Nodes don't agree about configuration!"
        else
            xputs "[OK] All nodes agree about slots configuration."
        end
    end

    def is_config_consistent?
        signatures=[]
        @nodes.each{|n|
            signatures << n.get_config_signature
        }
        return signatures.uniq.length == 1
    end

    def wait_cluster_join
        print "Waiting for the cluster to join"
        while !is_config_consistent?
            print "."
            STDOUT.flush
            sleep 1
        end
        print "\n"
    end

    def alloc_slots
        nodes_count = @nodes.length
        masters_count = @nodes.length / (@replicas+1)
        slots_per_node = ClusterHashSlots / masters_count
        masters = []
        slaves = []

        # The first step is to split instances by IP. This is useful as
        # we'll try to allocate master nodes in different physical machines
        # (as much as possible) and to allocate slaves of a given master in
        # different physical machines as well.
        #
        # This code assumes just that if the IP is different, than it is more
        # likely that the instance is running in a different physical host
        # or at least a different virtual machine.
        ips = {}
        @nodes.each{|n|
            ips[n.info[:host]] = [] if !ips[n.info[:host]]
            ips[n.info[:host]] << n
        }

        # Select master instances
        puts "Using #{masters_count} masters:"
        while masters.length < masters_count
            ips.each{|ip,nodes_list|
                next if nodes_list.length == 0
                masters << nodes_list.shift
                puts masters[-1]
                nodes_count -= 1
                break if masters.length == masters_count
            }
        end

        # Alloc slots on masters
        i = 0
        masters.each_with_index{|n,masternum|
            first = i*slots_per_node
            last = first+slots_per_node-1
            last = ClusterHashSlots-1 if masternum == masters.length-1
            n.add_slots first..last
            i += 1
        }

        # Select N replicas for every master.
        # We try to split the replicas among all the IPs with spare nodes
        # trying to avoid the host where the master is running, if possible.
        masters.each{|m|
            i = 0
            while i < @replicas
                ips.each{|ip,nodes_list|
                    next if nodes_list.length == 0
                    # Skip instances with the same IP as the master if we
                    # have some more IPs available.
                    next if ip == m.info[:host] && nodes_count > nodes_list.length
                    slave = nodes_list.shift
                    slave.set_as_replica(m.info[:name])
                    nodes_count -= 1
                    i += 1
                    puts "#{m} replica ##{i} is #{slave}"
                    break if masters.length == masters_count
                }
            end
        }
    end

    def flush_nodes_config
        @nodes.each{|n|
            n.flush_node_config
        }
    end

    def show_nodes
        @nodes.each{|n|
            xputs n.info_string
        }
    end

    def join_cluster
        # We use a brute force approach to make sure the node will meet
        # each other, that is, sending CLUSTER MEET messages to all the nodes
        # about the very same node.
        # Thanks to gossip this information should propagate across all the
        # cluster in a matter of seconds.
        first = false
        @nodes.each{|n|
            if !first then first = n.info; next; end # Skip the first node
            n.r.cluster("meet",first[:host],first[:port])
        }
    end

    def yes_or_die(msg)
        print "#{msg} (type 'yes' to accept): "
        STDOUT.flush
        if !(STDIN.gets.chomp.downcase == "yes")
            xputs "*** Aborting..."
            exit 1
        end
    end

    def load_cluster_info_from_node(nodeaddr)
        node = ClusterNode.new(nodeaddr)
        node.connect(:abort => true)
        node.assert_cluster
        node.load_info(:getfriends => true)
        add_node(node)
        node.friends.each{|f|
            next if f[:flags].index("noaddr") ||
                    f[:flags].index("disconnected") ||
                    f[:flags].index("fail")
            fnode = ClusterNode.new(f[:addr])
            fnode.connect()
            fnode.load_info()
            add_node(fnode)
        }
    end

    # Given a list of source nodes return a "resharding plan"
    # with what slots to move in order to move "numslots" slots to another
    # instance.
    def compute_reshard_table(sources,numslots)
        moved = []
        # Sort from bigger to smaller instance, for two reasons:
        # 1) If we take less slots than instances it is better to start
        #    getting from the biggest instances.
        # 2) We take one slot more from the first instance in the case of not
        #    perfect divisibility. Like we have 3 nodes and need to get 10
        #    slots, we take 4 from the first, and 3 from the rest. So the
        #    biggest is always the first.
        sources = sources.sort{|a,b| b.slots.length <=> a.slots.length}
        source_tot_slots = sources.inject(0) {|sum,source|
            sum+source.slots.length
        }
        sources.each_with_index{|s,i|
            # Every node will provide a number of slots proportional to the
            # slots it has assigned.
            n = (numslots.to_f/source_tot_slots*s.slots.length)
            if i == 0
                n = n.ceil
            else
                n = n.floor
            end
            s.slots.keys.sort[(0...n)].each{|slot|
                if moved.length < numslots
                    moved << {:source => s, :slot => slot}
                end
            }
        }
        return moved
    end

    def show_reshard_table(table)
        table.each{|e|
            puts "    Moving slot #{e[:slot]} from #{e[:source].info[:name]}"
        }
    end

    def move_slot(source,target,slot,o={})
        # We start marking the slot as importing in the destination node,
        # and the slot as migrating in the target host. Note that the order of
        # the operations is important, as otherwise a client may be redirected
        # to the target node that does not yet know it is importing this slot.
        print "Moving slot #{slot} from #{source} to #{target}: "; STDOUT.flush
        target.r.cluster("setslot",slot,"importing",source.info[:name])
        source.r.cluster("setslot",slot,"migrating",target.info[:name])
        # Migrate all the keys from source to target using the MIGRATE command
        while true
            keys = source.r.cluster("getkeysinslot",slot,10)
            break if keys.length == 0
            keys.each{|key|
                source.r.migrate(target.info[:host],target.info[:port],key,0,1000)
                print "." if o[:verbose]
                STDOUT.flush
            }
        end
        puts
        # Set the new node as the owner of the slot in all the known nodes.
        @nodes.each{|n|
            n.r.cluster("setslot",slot,"node",target.info[:name])
        }
    end

    # redis-trib subcommands implementations

    def check_cluster_cmd(argv,opt)
        load_cluster_info_from_node(argv[0])
        check_cluster
    end

    def fix_cluster_cmd(argv,opt)
        @fix = true
        load_cluster_info_from_node(argv[0])
        check_cluster
    end

    def reshard_cluster_cmd(argv,opt)
        load_cluster_info_from_node(argv[0])
        check_cluster
        if @errors.length != 0
            puts "*** Please fix your cluster problems before resharding"
            exit 1
        end
        numslots = 0
        while numslots <= 0 or numslots > ClusterHashSlots
            print "How many slots do you want to move (from 1 to #{ClusterHashSlots})? "
            numslots = STDIN.gets.to_i
        end
        target = nil
        while not target
            print "What is the receiving node ID? "
            target = get_node_by_name(STDIN.gets.chop)
            if !target || target.has_flag?("slave")
                xputs "*** The specified node is not known or not a master, please retry."
                target = nil
            end
        end
        sources = []
        xputs "Please enter all the source node IDs."
        xputs "  Type 'all' to use all the nodes as source nodes for the hash slots."
        xputs "  Type 'done' once you entered all the source nodes IDs."
        while true
            print "Source node ##{sources.length+1}:"
            line = STDIN.gets.chop
            src = get_node_by_name(line)
            if line == "done"
                if sources.length == 0
                    puts "No source nodes given, operation aborted"
                    exit 1
                else
                    break
                end
            elsif line == "all"
                @nodes.each{|n|
                    next if n.info[:name] == target.info[:name]
                    next if n.has_flag?("slave")
                    sources << n
                }
                break
            elsif !src || src.has_flag?("slave")
                xputs "*** The specified node is not known or is not a master, please retry."
            elsif src.info[:name] == target.info[:name]
                xputs "*** It is not possible to use the target node as source node."
            else
                sources << src
            end
        end
        puts "\nReady to move #{numslots} slots."
        puts "  Source nodes:"
        sources.each{|s| puts "    "+s.info_string}
        puts "  Destination node:"
        puts "    #{target.info_string}"
        reshard_table = compute_reshard_table(sources,numslots)
        puts "  Resharding plan:"
        show_reshard_table(reshard_table)
        print "Do you want to proceed with the proposed reshard plan (yes/no)? "
        yesno = STDIN.gets.chop
        exit(1) if (yesno != "yes")
        reshard_table.each{|e|
            move_slot(e[:source],target,e[:slot],:verbose=>true)
        }
    end

    # This is an helper function for create_cluster_cmd that verifies if
    # the number of nodes and the specified replicas have a valid configuration
    # where there are at least three master nodes and enough replicas per node.
    def check_create_parameters
        masters = @nodes.length/(@replicas+1)
        if masters < 3
            puts "*** ERROR: Invalid configuration for cluster creation."
            puts "*** Redis Cluster requires at least 3 master nodes."
            puts "*** This is not possible with #{@nodes.length} nodes and #{@replicas} replicas per node."
            puts "*** At least #{3*(@replicas+1)} nodes are required."
            exit 1
        end
    end

    def create_cluster_cmd(argv,opt)
        opt = {'replicas' => 0}.merge(opt)
        @replicas = opt['replicas'].to_i

        xputs ">>> Creating cluster"
        argv[0..-1].each{|n|
            node = ClusterNode.new(n)
            node.connect(:abort => true)
            node.assert_cluster
            node.load_info
            node.assert_empty
            add_node(node)
        }
        check_create_parameters
        xputs ">>> Performing hash slots allocation on #{@nodes.length} nodes..."
        alloc_slots
        show_nodes
        yes_or_die "Can I set the above configuration?"
        flush_nodes_config
        xputs ">>> Nodes configuration updated"
        xputs ">>> Sending CLUSTER MEET messages to join the cluster"
        join_cluster
        # Give one second for the join to start, in order to avoid that
        # wait_cluster_join will find all the nodes agree about the config as
        # they are still empty with unassigned slots.
        sleep 1
        wait_cluster_join
        flush_nodes_config # Useful for the replicas
        check_cluster
    end

    def addnode_cluster_cmd(argv,opt)
        xputs ">>> Adding node #{argv[0]} to cluster #{argv[1]}"

        # Check the existing cluster
        load_cluster_info_from_node(argv[1])
        check_cluster

        # Add the new node
        new = ClusterNode.new(argv[0])
        new.connect(:abort => true)
        new.assert_cluster
        new.load_info
        new.assert_empty
        first = @nodes.first.info

        # Send CLUSTER MEET command to the new node
        xputs ">>> Send CLUSTER MEET to node #{new} to make it join the cluster."
        new.r.cluster("meet",first[:host],first[:port])
    end

    def help_cluster_cmd(opt)
        show_help
        exit 0
    end

    # Parse the options for the specific command "cmd".
    # Returns an hash populate with option => value pairs, and the index of
    # the first non-option argument in ARGV.
    def parse_options(cmd)
        idx = 1 ; # Current index into ARGV
        options={}
        while idx < ARGV.length && ARGV[idx][0..1] == '--'
            if ARGV[idx][0..1] == "--"
                option = ARGV[idx][2..-1]
                idx += 1
                if ALLOWED_OPTIONS[cmd] == nil || ALLOWED_OPTIONS[cmd][option] == nil
                    puts "Unknown option '#{option}' for command '#{cmd}'"
                    exit 1
                end
                if ALLOWED_OPTIONS[cmd][option]
                    value = ARGV[idx]
                    idx += 1
                else
                    value = true
                end
                options[option] = value
            else
                # Remaining arguments are not options.
                break
            end
        end
        return options,idx
    end
end

COMMANDS={
    "create"  => ["create_cluster_cmd", -2, "host1:port1 ... hostN:portN"],
    "check"   => ["check_cluster_cmd", 2, "host:port"],
    "fix"     => ["fix_cluster_cmd", 2, "host:port"],
    "reshard" => ["reshard_cluster_cmd", 2, "host:port"],
    "addnode" => ["addnode_cluster_cmd", 3, "new_host:new_port existing_host:existing_port"],
    "help"    => ["help_cluster_cmd", 1, "(show this help)"]
}

ALLOWED_OPTIONS={
    "create" => {"replicas" => true}
}

def show_help
    puts "Usage: redis-trib <command> <options> <arguments ...>"
    puts
    COMMANDS.each{|k,v|
        puts "  #{k.ljust(10)} #{v[2]}"
    }
    puts
end

# Sanity check
if ARGV.length == 0
    show_help
    exit 1
end

rt = RedisTrib.new
cmd_spec = COMMANDS[ARGV[0].downcase]
if !cmd_spec
    puts "Unknown redis-trib subcommand '#{ARGV[0]}'"
    exit 1
end

# Parse options
cmd_options,first_non_option = rt.parse_options(ARGV[0].downcase)
rt.check_arity(cmd_spec[1],ARGV.length-(first_non_option-1))

# Dispatch
rt.send(cmd_spec[0],ARGV[first_non_option..-1],cmd_options)
