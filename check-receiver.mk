datasource_programs += [
  ( "cat /var/lib/icinga/ramdisk/check-receiver.<HOST>", [ 'nagpush' ], ALL_HOSTS ),
]

# Add host check for stale data files
extra_nagios_conf += r"""
define command {
    command_name check-receiver-hostcheck
    command_line $USER1$/check_file_age -w 60 -c 90 -f /var/lib/icinga/ramdisk/check-receiver.$HOSTNAME$
}
"""

# use the above defined host check
extra_host_conf["check_command"] = [("check-receiver-hostcheck", [ "nagpush" ], ALL_HOSTS )]
