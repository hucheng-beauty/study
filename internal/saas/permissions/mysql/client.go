package mysql

import (
    "context"
    "fmt"

    "study/internal/saas/permissions"

    "github.com/jmoiron/sqlx"

    // "code.byted.org/lab-speech/saas_backend/auth"
    // "code.byted.org/lab-speech/saas_backend/permissions"
)

type ProjectPermission struct {
    Username       string `db:"username"`
    Name           string `db:"name"`
    Email          string `db:"email"`
    PermissionType string `db:"type"`
}

type ProjectPermissions []ProjectPermission

func (rp ProjectPermissions) PermissionHolders(parcelID int) permissions.PermissionHolders {
    var permissionHolders permissions.PermissionHolders
projectPermissionLoop:
    for _, projectPermission := range rp {
        permissionDescriptor := permissions.PermissionType(projectPermission.PermissionType).Descriptor(parcelID)
        temporality := "temporary"
        if projectPermission.PermissionType == "admin" {
            temporality = "permanent"
        }
        permission := permissions.Permission{
            Descriptor:  permissionDescriptor,
            Temporality: temporality,
        }
        // When we add external user to database via permissions adding
        // he has no name and no email inserted.
        // We change his name and email with username in order to display something
        // on front-end page.
        name := projectPermission.Name
        if name == "" {
            name = projectPermission.Username
        }
        email := projectPermission.Email
        if email == "" {
            email = projectPermission.Username + "@bytedance.com"
        }
        for i := range permissionHolders {
            if permissionHolders[i].Email == email {
                permissionHolders[i].Permissions = append(permissionHolders[i].Permissions, permission)
                continue projectPermissionLoop
            }
        }
        permissionHolders = append(permissionHolders, permissions.PermissionHolder{
            User: permissions.User{
                Name:  name,
                Email: email,
            },
            Target:      "user",
            Permissions: permissions.Permissions{permission},
        })
    }
    return permissionHolders
}

type dbPermissionsClient struct {
    db *sqlx.DB
}

func NewDBPermissionsClient(db *sqlx.DB) DBPermissionsClient {
    return dbPermissionsClient{db: db}
}

type UserAccessControl struct {
    ProjectID  int    `db:"project_id"`
    Permission string `db:"permission_type"`
}

type UserAccess []UserAccessControl

func (ra UserAccess) PermissionsUserAccess() permissions.UserAccess {
    var permissionsUserAccess permissions.UserAccess
accessLoop:
    for i := range ra {
        for j := range permissionsUserAccess {
            if ra[i].ProjectID == permissionsUserAccess[j].ProjectID {
                permissionsUserAccess[j].UserPermissions = append(permissionsUserAccess[j].UserPermissions,
                    permissions.PermissionType(ra[i].Permission))
                continue accessLoop
            }
        }
        permissionsUserAccess = append(permissionsUserAccess, permissions.UserAccessControl{
            ProjectID:       ra[i].ProjectID,
            UserPermissions: []permissions.PermissionType{permissions.PermissionType(ra[i].Permission)},
        })
    }
    return permissionsUserAccess
}

func (c dbPermissionsClient) UserPermittedProjects(ctx context.Context, user auth.User) (UserAccess, error) {
    var userAccess UserAccess
    err := c.db.SelectContext(ctx, &userAccess, `
		SELECT 
			project_id, type as permission_type
		FROM
			permission
		JOIN
			user on user.id=user_id
		WHERE
			username=?
		ORDER BY
			project_id;`, user.Username)
    // ORDER BY only for deterministic result needed for testing purposes.
    if err != nil {
        return nil, fmt.Errorf("select user allowed project ids: %w", err)
    }
    return userAccess, nil
}

func (c dbPermissionsClient) ProjectPermissions(ctx context.Context,
        projectID int) (ProjectPermissions, error) {
    var projectPermissions ProjectPermissions
    err := c.db.SelectContext(ctx, &projectPermissions, `
		SELECT
			u.username, u.name, u.email, p.type
		FROM
			permission p
		JOIN 
			user u on u.id = p.user_id
		WHERE
			project_id = ?
		ORDER BY
			u.name;`, projectID)
    // ORDER BY only for deterministic result needed for testing purposes.
    if err != nil {
        return nil, fmt.Errorf("select project users with permissions: %w", err)
    }
    return projectPermissions, nil
}

func (c dbPermissionsClient) DeletePermission(ctx context.Context, projectID int,
        user auth.User, perm permissions.PermissionType) error {
    _, err := c.db.ExecContext(ctx, `
		DELETE FROM
			permission
		WHERE
			project_id = ? 
		    and user_id = (SELECT id from user where username=?) 
		    and type = ?;`, projectID, user.Username, string(perm))
    if err != nil {
        return fmt.Errorf("delete permission %s from user %s to project %d: %w",
            string(perm), user.Username, projectID, err)
    }
    return nil
}

func (c dbPermissionsClient) AddPermission(ctx context.Context, projectID int,
        user auth.User, perm permissions.PermissionType) error {
    err := c.ensureUpToDateUserExistence(ctx, user)
    if err != nil {
        return fmt.Errorf("ensure user %s existence: %w", user.Username, err)
    }
    _, err = c.db.ExecContext(ctx, `
		INSERT INTO
			permission (project_id, user_id, type)
		VALUES 
			(?, (SELECT id from user where username=?), ?)
		ON DUPLICATE KEY UPDATE type=type;`,
        projectID, user.Username, string(perm))
    if err != nil {
        return fmt.Errorf("add permission %s to user %s for project %d: %w",
            string(perm), user.Username, projectID, err)
    }
    return nil
}

func (c dbPermissionsClient) ensureUpToDateUserExistence(ctx context.Context, user auth.User) error {
    var err error
    if user.Name == "" || user.Email == "" {
        // If user has missing fields, we don't update values stored in Data Base
        // since we may lose information about user.
        _, err = c.db.ExecContext(ctx, `
			INSERT INTO
				user (username, name, email)
			VALUES
				(?, ?, ?)
			ON DUPLICATE KEY UPDATE name=name, email=email;`,
            user.Username, user.Name, user.Email)
    } else {
        // If user has all fields initialized, we perform update, since it may change "unknown user"
        // aka only-username-present-record into fully-known user aka all-fields-fulfilled-record.
        _, err = c.db.ExecContext(ctx, `
			INSERT INTO
				user (username, name, email)
			VALUES
				(?, ?, ?)
			ON DUPLICATE KEY UPDATE name=?, email=?;`,
            user.Username, user.Name, user.Email, user.Name, user.Email)
    }
    if err != nil {
        return fmt.Errorf("ensure user %s existence in database: %w", user.Username, err)
    }
    return nil
}
