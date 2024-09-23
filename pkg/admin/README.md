# FreePDM
A PDM for FreeCAD

## The Admin page. 

This page is a long list of things that need to be done. Don't expect everything working. It just isn't done yet. The list of TODO's is very long.

### User accounts
- Login (/admin/login) - a simple login form with username admin and password admin.
- Rename (/admin/rename) - a form to rename the admin account.
- Home (/admin) - the home page for the admin user.
- User List (/admin/users) - the default page of the users section for performing CRUD operations. Displays a list of all users with buttons to add (/admin/user/user_name/add) and edit (/admin/user/user_name).
- User (/admin/user/user_name) - a form to modify a specific user acccount.
- Unactivate (/admin/user/user_name/unactivate) and reactivate (/admin/user/user_name/reactivate).

### Vaults
- Vaults (/admin/vaults) - a list of the existing vaults, incl. buttons for add, edit, remove.
- Vault (/admin/vault/vault_name) Setting up name of the vault, it's versioning scheme, numbering scheme and wether it is mandatory. Also the user access. Who has access and what role?
- Unlock (/admin/vault/vault_name/unlock) a file or item revision.
- Unremove (/admin/vault/vault_name/unremove) "Un-remove" a file or an item-revision.
- Vault logs (/admin/vault/logs) Show the logs.

- And a lot more...