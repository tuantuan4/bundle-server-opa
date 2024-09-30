
package app.rbac
import rego.v1

# By default, deny requests.
default allow := false

# Allow the action if the user is granted permission to perform the action.
allow if {
	# Find grants for the user.
	some grant in user_is_granted

	# Check if the grant permits the action.
	input.method == grant.method
	input.url == grant.url
}


# user_is_granted is a set of grants for the user identified in the request.
# The `grant` will be contained if the set `user_is_granted` for every...
user_is_granted contains grant if {
	# `role` assigned an element of the user_roles for this user...
	some role in data.bundle.data.user_roles[input.username]

	# `grant` assigned a single grant from the grants list for 'role'...
	some grant in data.bundle.data.role_grants[role]
}
