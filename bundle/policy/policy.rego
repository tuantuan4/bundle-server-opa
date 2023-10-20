package main

default allow = false

# Kiểm tra xem một người dùng có quyền là "admin" hay không
allow {
    user = input[i]
    user.role == "admin"
}

