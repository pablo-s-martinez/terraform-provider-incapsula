---
layout: "incapsula"
page_title: "Incapsula: user"
sidebar_current: "docs-incapsula-resource-user"
description: |-
  Provides a Incapsula User resource.
---

# incapsula_site

Provides a Incapsula SubAccount resource. 

## Example Usage

```hcl
resource "incapsula_user" "example-user" {
  email             = "test@test.com"
  account_id        = "1190270"
  role_names        = ["Reader"]
  first_name        = "first"
  last_name         = "last"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) Email address. For example: joe@example.com. example: `userEmail@imperva.com`.
* `account_id` - (Required) Unique ID of the required account . example: 123456.
* `first_name` - (Optional) The first name of the user that was acted on. example: `John`.
* `last_name` - (Optional) The last name of the user that was acted on. example: `Snow`.
* `role_names` - (Optional) List of role names to add to the user. for example : [`Reader`].

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the User.

