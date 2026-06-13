-- 000005_ipam.down.sql
DROP TABLE IF EXISTS dns_records;
DROP TABLE IF EXISTS dns_zones;
DROP TABLE IF EXISTS dhcp_leases;
DROP TABLE IF EXISTS ip_addresses;
DROP TABLE IF EXISTS prefixes;
ALTER TABLE network_interfaces DROP CONSTRAINT IF EXISTS fk_ni_untagged_vlan;
DROP TABLE IF EXISTS vlans;
DROP TABLE IF EXISTS vlan_groups;
DROP TABLE IF EXISTS vrfs;
