# FreePDM
A PDM for FreeCAD.

After a long time of thinking I came to the conclusion that only having JSON interfaces with the client would be sufficient. That is, I think, *now* is the best way to communicate with the client. I don't think that having a web page is the right choice, but having a client is a better approach because then the server would have to do a lot less, which means more stable, and also the client can be anything such as a Windows app, Linux app and a Mac app, and maybe in the long future a mobile app.

## The PDM Server page

This page is a long list of things that need to be done. Don't expect one thing working. It just isn't done yet. The list of TODO's is very long.

### Basic functionality
- Commands (/command) for ordinary functions such a ls, mv, rename, copy, chdir

### User accounts
- Login (/admin/login) - a simple login form with username admin and password admin.
- Rename (/admin/rename) - a form to rename the admin account.
- Home (/admin) - the home page for the admin user.
- User List (/admin/users) - the default page of the users section for performing CRUD operations. Displays a list of all users with buttons to add (/admin/user/user_name/add) and edit (/admin/user/user_name).
- User (/admin/user/user_name) - a form to modify a specific user account.
- Unactivate (/admin/user/user_name/unactivate) and reactivate (/admin/user/user_name/reactivate).
- Vault (/admin/vault/vault_name) Setting up name of the vault, it's versioning scheme, numbering scheme and whether it is mandatory. Also the user access. 
- Unremove (/admin/vault/vault_name/unremove) "Un-remove" a file or an item-revision.
- Vault logs (/admin/vault/logs) Show the logs.

### Vaults
- Vaults (/vaults/list) - a list of the existing vaults.
Who has access and what role?
- Unlock (/vault/vault_name/unlock) a file.

- And a lot more...