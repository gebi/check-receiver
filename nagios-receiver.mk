datasource_programs += [
  ( "cat /var/lib/icinga/ramdisk/nagios-receiver.<HOST>", [ 'nagpush' ], ALL_HOSTS ),
]

# Add host check for stale data files
extra_nagios_conf += r"""
define command {
    command_name nagios-receiver-hostcheck
    command_line $USER1$/check_file_age -w 60 -c 90 -f /var/lib/icinga/ramdisk/nagios-receiver.$HOSTNAME$
}
"""

# use the above defined host check
extra_host_conf["check_command"] = [("nagios-receiver-hostcheck", [ "nagpush" ], ALL_HOSTS )]
