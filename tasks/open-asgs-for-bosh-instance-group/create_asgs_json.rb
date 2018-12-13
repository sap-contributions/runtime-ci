#!/usr/bin/ruby -w
# frozen_string_literal: true

require 'json'
require 'open3'

def write_asgs_json(instance_group_ips, destination_file)
  asgs = []

  instance_group_ips.each do |ip|
    asgs << { protocol: 'tcp', destination: ip, ports: '1-65535' }
  end

  File.open(destination_file, 'w') do |f|
    f.write(asgs.to_json)
  end
end

def get_ips_from_bosh_output(instance_group_name)
  instance_ips = []

  stdout, _, exitcode = Open3.capture3('bosh is --json')

  raise "'bosh is --json' returned an error: #{stdout}" if exitcode != 0

  bosh_tables = JSON.parse(stdout)['Tables']
  raise 'More than one bosh deployment detected' unless bosh_tables.size == 1

  instances = bosh_tables.first['Rows']
  instances.each do |is|
    if is['instance'].include? instance_group_name
      instance_ips << is['ips'].split(/\s/).select { |ip| ip.start_with? '10.' }
    end
  end

  instance_ips.flatten!

  raise 'No IPs detected' if instance_ips.empty?

  instance_ips
end

instance_name, destination_file = ARGV

instance_group_ips = get_ips_from_bosh_output(instance_name)
write_asgs_json(instance_group_ips, destination_file)
