package app.rbac

import rego.v1

default allow := false

allow if {
	some role in data.bundle.data.user_roles[input.username]

	data.bundle.data.role_grants[role][input.method][input.url]
}

