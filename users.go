//
// Copyright 2017, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"errors"
	"fmt"
	"time"
)

// UsersService handles communication with the user related methods of
// the GitLab API.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html
type UsersService struct {
	client *Client
}

// User represents a GitLab user.
//
// GitLab API docs: https://docs.gitlab.com/ee/api/users.html
type User struct {
	ID                        int                `bson:"id" json:"id"`
	Username                  string             `bson:"username" json:"username"`
	Email                     string             `bson:"email" json:"email"`
	Name                      string             `bson:"name" json:"name"`
	State                     string             `bson:"state" json:"state"`
	CreatedAt                 *time.Time         `bson:"created_at" json:"created_at"`
	Bio                       string             `bson:"bio" json:"bio"`
	Location                  string             `bson:"location" json:"location"`
	PublicEmail               string             `bson:"public_email" json:"public_email"`
	Skype                     string             `bson:"skype" json:"skype"`
	Linkedin                  string             `bson:"linkedin" json:"linkedin"`
	Twitter                   string             `bson:"twitter" json:"twitter"`
	WebsiteURL                string             `bson:"website_url" json:"website_url"`
	Organization              string             `bson:"organization" json:"organization"`
	ExternUID                 string             `bson:"extern_uid" json:"extern_uid"`
	Provider                  string             `bson:"provider" json:"provider"`
	ThemeID                   int                `bson:"theme_id" json:"theme_id"`
	LastActivityOn            *ISOTime           `bson:"last_activity_on" json:"last_activity_on"`
	ColorSchemeID             int                `bson:"color_scheme_id" json:"color_scheme_id"`
	IsAdmin                   bool               `bson:"is_admin" json:"is_admin"`
	AvatarURL                 string             `bson:"avatar_url" json:"avatar_url"`
	CanCreateGroup            bool               `bson:"can_create_group" json:"can_create_group"`
	CanCreateProject          bool               `bson:"can_create_project" json:"can_create_project"`
	ProjectsLimit             int                `bson:"projects_limit" json:"projects_limit"`
	CurrentSignInAt           *time.Time         `bson:"current_sign_in_at" json:"current_sign_in_at"`
	LastSignInAt              *time.Time         `bson:"last_sign_in_at" json:"last_sign_in_at"`
	ConfirmedAt               *time.Time         `bson:"confirmed_at" json:"confirmed_at"`
	TwoFactorEnabled          bool               `bson:"two_factor_enabled" json:"two_factor_enabled"`
	Identities                []*UserIdentity    `bson:"identities" json:"identities"`
	External                  bool               `bson:"external" json:"external"`
	PrivateProfile            bool               `bson:"private_profile" json:"private_profile"`
	SharedRunnersMinutesLimit int                `bson:"shared_runners_minutes_limit" json:"shared_runners_minutes_limit"`
	CustomAttributes          []*CustomAttribute `bson:"custom_attributes" json:"custom_attributes"`
}

// UserIdentity represents a user identity.
type UserIdentity struct {
	Provider  string `bson:"provider" json:"provider"`
	ExternUID string `bson:"extern_uid" json:"extern_uid"`
}

// ListUsersOptions represents the available ListUsers() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#list-users
type ListUsersOptions struct {
	ListOptions
	Active  *bool `url:"active,omitempty" bson:"active,omitempty" json:"active,omitempty"`
	Blocked *bool `url:"blocked,omitempty" bson:"blocked,omitempty" json:"blocked,omitempty"`

	// The options below are only available for admins.
	Search               *string    `url:"search,omitempty" bson:"search,omitempty" json:"search,omitempty"`
	Username             *string    `url:"username,omitempty" bson:"username,omitempty" json:"username,omitempty"`
	ExternalUID          *string    `url:"extern_uid,omitempty" bson:"extern_uid,omitempty" json:"extern_uid,omitempty"`
	Provider             *string    `url:"provider,omitempty" bson:"provider,omitempty" json:"provider,omitempty"`
	CreatedBefore        *time.Time `url:"created_before,omitempty" bson:"created_before,omitempty" json:"created_before,omitempty"`
	CreatedAfter         *time.Time `url:"created_after,omitempty" bson:"created_after,omitempty" json:"created_after,omitempty"`
	OrderBy              *string    `url:"order_by,omitempty" bson:"order_by,omitempty" json:"order_by,omitempty"`
	Sort                 *string    `url:"sort,omitempty" bson:"sort,omitempty" json:"sort,omitempty"`
	WithCustomAttributes *bool      `url:"with_custom_attributes,omitempty" bson:"with_custom_attributes,omitempty" json:"with_custom_attributes,omitempty"`
}

// ListUsers gets a list of users.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#list-users
func (s *UsersService) ListUsers(opt *ListUsersOptions, options ...OptionFunc) ([]*User, *Response, error) {
	req, err := s.client.NewRequest("GET", "users", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var usr []*User
	resp, err := s.client.Do(req, &usr)
	if err != nil {
		return nil, resp, err
	}

	return usr, resp, err
}

// GetUser gets a single user.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#single-user
func (s *UsersService) GetUser(user int, options ...OptionFunc) (*User, *Response, error) {
	u := fmt.Sprintf("users/%d", user)

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}

	return usr, resp, err
}

// CreateUserOptions represents the available CreateUser() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#user-creation
type CreateUserOptions struct {
	Email            *string `url:"email,omitempty" bson:"email,omitempty" json:"email,omitempty"`
	Password         *string `url:"password,omitempty" bson:"password,omitempty" json:"password,omitempty"`
	ResetPassword    *bool   `url:"reset_password,omitempty" bson:"reset_password,omitempty" json:"reset_password,omitempty"`
	Username         *string `url:"username,omitempty" bson:"username,omitempty" json:"username,omitempty"`
	Name             *string `url:"name,omitempty" bson:"name,omitempty" json:"name,omitempty"`
	Skype            *string `url:"skype,omitempty" bson:"skype,omitempty" json:"skype,omitempty"`
	Linkedin         *string `url:"linkedin,omitempty" bson:"linkedin,omitempty" json:"linkedin,omitempty"`
	Twitter          *string `url:"twitter,omitempty" bson:"twitter,omitempty" json:"twitter,omitempty"`
	WebsiteURL       *string `url:"website_url,omitempty" bson:"website_url,omitempty" json:"website_url,omitempty"`
	Organization     *string `url:"organization,omitempty" bson:"organization,omitempty" json:"organization,omitempty"`
	ProjectsLimit    *int    `url:"projects_limit,omitempty" bson:"projects_limit,omitempty" json:"projects_limit,omitempty"`
	ExternUID        *string `url:"extern_uid,omitempty" bson:"extern_uid,omitempty" json:"extern_uid,omitempty"`
	Provider         *string `url:"provider,omitempty" bson:"provider,omitempty" json:"provider,omitempty"`
	Bio              *string `url:"bio,omitempty" bson:"bio,omitempty" json:"bio,omitempty"`
	Location         *string `url:"location,omitempty" bson:"location,omitempty" json:"location,omitempty"`
	Admin            *bool   `url:"admin,omitempty" bson:"admin,omitempty" json:"admin,omitempty"`
	CanCreateGroup   *bool   `url:"can_create_group,omitempty" bson:"can_create_group,omitempty" json:"can_create_group,omitempty"`
	SkipConfirmation *bool   `url:"skip_confirmation,omitempty" bson:"skip_confirmation,omitempty" json:"skip_confirmation,omitempty"`
	External         *bool   `url:"external,omitempty" bson:"external,omitempty" json:"external,omitempty"`
}

// CreateUser creates a new user. Note only administrators can create new users.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#user-creation
func (s *UsersService) CreateUser(opt *CreateUserOptions, options ...OptionFunc) (*User, *Response, error) {
	req, err := s.client.NewRequest("POST", "users", opt, options)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}

	return usr, resp, err
}

// ModifyUserOptions represents the available ModifyUser() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#user-modification
type ModifyUserOptions struct {
	Email              *string `url:"email,omitempty" bson:"email,omitempty" json:"email,omitempty"`
	Password           *string `url:"password,omitempty" bson:"password,omitempty" json:"password,omitempty"`
	Username           *string `url:"username,omitempty" bson:"username,omitempty" json:"username,omitempty"`
	Name               *string `url:"name,omitempty" bson:"name,omitempty" json:"name,omitempty"`
	Skype              *string `url:"skype,omitempty" bson:"skype,omitempty" json:"skype,omitempty"`
	Linkedin           *string `url:"linkedin,omitempty" bson:"linkedin,omitempty" json:"linkedin,omitempty"`
	Twitter            *string `url:"twitter,omitempty" bson:"twitter,omitempty" json:"twitter,omitempty"`
	WebsiteURL         *string `url:"website_url,omitempty" bson:"website_url,omitempty" json:"website_url,omitempty"`
	Organization       *string `url:"organization,omitempty" bson:"organization,omitempty" json:"organization,omitempty"`
	ProjectsLimit      *int    `url:"projects_limit,omitempty" bson:"projects_limit,omitempty" json:"projects_limit,omitempty"`
	ExternUID          *string `url:"extern_uid,omitempty" bson:"extern_uid,omitempty" json:"extern_uid,omitempty"`
	Provider           *string `url:"provider,omitempty" bson:"provider,omitempty" json:"provider,omitempty"`
	Bio                *string `url:"bio,omitempty" bson:"bio,omitempty" json:"bio,omitempty"`
	Location           *string `url:"location,omitempty" bson:"location,omitempty" json:"location,omitempty"`
	Admin              *bool   `url:"admin,omitempty" bson:"admin,omitempty" json:"admin,omitempty"`
	CanCreateGroup     *bool   `url:"can_create_group,omitempty" bson:"can_create_group,omitempty" json:"can_create_group,omitempty"`
	SkipReconfirmation *bool   `url:"skip_reconfirmation,omitempty" bson:"skip_reconfirmation,omitempty" json:"skip_reconfirmation,omitempty"`
	External           *bool   `url:"external,omitempty" bson:"external,omitempty" json:"external,omitempty"`
}

// ModifyUser modifies an existing user. Only administrators can change attributes
// of a user.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#user-modification
func (s *UsersService) ModifyUser(user int, opt *ModifyUserOptions, options ...OptionFunc) (*User, *Response, error) {
	u := fmt.Sprintf("users/%d", user)

	req, err := s.client.NewRequest("PUT", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}

	return usr, resp, err
}

// DeleteUser deletes a user. Available only for administrators. This is an
// idempotent function, calling this function for a non-existent user id still
// returns a status code 200 OK. The JSON response differs if the user was
// actually deleted or not. In the former the user is returned and in the
// latter not.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#user-deletion
func (s *UsersService) DeleteUser(user int, options ...OptionFunc) (*Response, error) {
	u := fmt.Sprintf("users/%d", user)

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// CurrentUser gets currently authenticated user.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#current-user
func (s *UsersService) CurrentUser(options ...OptionFunc) (*User, *Response, error) {
	req, err := s.client.NewRequest("GET", "user", nil, options)
	if err != nil {
		return nil, nil, err
	}

	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}

	return usr, resp, err
}

// SSHKey represents a SSH key.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#list-ssh-keys
type SSHKey struct {
	ID        int        `bson:"id" json:"id"`
	Title     string     `bson:"title" json:"title"`
	Key       string     `bson:"key" json:"key"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
}

// ListSSHKeys gets a list of currently authenticated user's SSH keys.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#list-ssh-keys
func (s *UsersService) ListSSHKeys(options ...OptionFunc) ([]*SSHKey, *Response, error) {
	req, err := s.client.NewRequest("GET", "user/keys", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var k []*SSHKey
	resp, err := s.client.Do(req, &k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, err
}

// ListSSHKeysForUserOptions represents the available ListSSHKeysForUser() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#list-ssh-keys-for-user
type ListSSHKeysForUserOptions ListOptions

// ListSSHKeysForUser gets a list of a specified user's SSH keys. Available
// only for admin
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#list-ssh-keys-for-user
func (s *UsersService) ListSSHKeysForUser(user int, opt *ListSSHKeysForUserOptions, options ...OptionFunc) ([]*SSHKey, *Response, error) {
	u := fmt.Sprintf("users/%d/keys", user)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var k []*SSHKey
	resp, err := s.client.Do(req, &k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, err
}

// GetSSHKey gets a single key.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#single-ssh-key
func (s *UsersService) GetSSHKey(key int, options ...OptionFunc) (*SSHKey, *Response, error) {
	u := fmt.Sprintf("user/keys/%d", key)

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	k := new(SSHKey)
	resp, err := s.client.Do(req, k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, err
}

// AddSSHKeyOptions represents the available AddSSHKey() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/projects.html#add-ssh-key
type AddSSHKeyOptions struct {
	Title *string `url:"title,omitempty" bson:"title,omitempty" json:"title,omitempty"`
	Key   *string `url:"key,omitempty" bson:"key,omitempty" json:"key,omitempty"`
}

// AddSSHKey creates a new key owned by the currently authenticated user.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#add-ssh-key
func (s *UsersService) AddSSHKey(opt *AddSSHKeyOptions, options ...OptionFunc) (*SSHKey, *Response, error) {
	req, err := s.client.NewRequest("POST", "user/keys", opt, options)
	if err != nil {
		return nil, nil, err
	}

	k := new(SSHKey)
	resp, err := s.client.Do(req, k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, err
}

// AddSSHKeyForUser creates new key owned by specified user. Available only for
// admin.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#add-ssh-key-for-user
func (s *UsersService) AddSSHKeyForUser(user int, opt *AddSSHKeyOptions, options ...OptionFunc) (*SSHKey, *Response, error) {
	u := fmt.Sprintf("users/%d/keys", user)

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	k := new(SSHKey)
	resp, err := s.client.Do(req, k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, err
}

// DeleteSSHKey deletes key owned by currently authenticated user. This is an
// idempotent function and calling it on a key that is already deleted or not
// available results in 200 OK.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#delete-ssh-key-for-current-owner
func (s *UsersService) DeleteSSHKey(key int, options ...OptionFunc) (*Response, error) {
	u := fmt.Sprintf("user/keys/%d", key)

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteSSHKeyForUser deletes key owned by a specified user. Available only
// for admin.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#delete-ssh-key-for-given-user
func (s *UsersService) DeleteSSHKeyForUser(user, key int, options ...OptionFunc) (*Response, error) {
	u := fmt.Sprintf("users/%d/keys/%d", user, key)

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// BlockUser blocks the specified user. Available only for admin.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#block-user
func (s *UsersService) BlockUser(user int, options ...OptionFunc) error {
	u := fmt.Sprintf("users/%d/block", user)

	req, err := s.client.NewRequest("POST", u, nil, options)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 201:
		return nil
	case 403:
		return errors.New("Cannot block a user that is already blocked by LDAP synchronization")
	case 404:
		return errors.New("User does not exist")
	default:
		return fmt.Errorf("Received unexpected result code: %d", resp.StatusCode)
	}
}

// UnblockUser unblocks the specified user. Available only for admin.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#unblock-user
func (s *UsersService) UnblockUser(user int, options ...OptionFunc) error {
	u := fmt.Sprintf("users/%d/unblock", user)

	req, err := s.client.NewRequest("POST", u, nil, options)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 201:
		return nil
	case 403:
		return errors.New("Cannot unblock a user that is blocked by LDAP synchronization")
	case 404:
		return errors.New("User does not exist")
	default:
		return fmt.Errorf("Received unexpected result code: %d", resp.StatusCode)
	}
}

// Email represents an Email.
//
// GitLab API docs: https://doc.gitlab.com/ce/api/users.html#list-emails
type Email struct {
	ID    int    `bson:"id" json:"id"`
	Email string `bson:"email" json:"email"`
}

// ListEmails gets a list of currently authenticated user's Emails.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#list-emails
func (s *UsersService) ListEmails(options ...OptionFunc) ([]*Email, *Response, error) {
	req, err := s.client.NewRequest("GET", "user/emails", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var e []*Email
	resp, err := s.client.Do(req, &e)
	if err != nil {
		return nil, resp, err
	}

	return e, resp, err
}

// ListEmailsForUserOptions represents the available ListEmailsForUser() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#list-emails-for-user
type ListEmailsForUserOptions ListOptions

// ListEmailsForUser gets a list of a specified user's Emails. Available
// only for admin
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#list-emails-for-user
func (s *UsersService) ListEmailsForUser(user int, opt *ListEmailsForUserOptions, options ...OptionFunc) ([]*Email, *Response, error) {
	u := fmt.Sprintf("users/%d/emails", user)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var e []*Email
	resp, err := s.client.Do(req, &e)
	if err != nil {
		return nil, resp, err
	}

	return e, resp, err
}

// GetEmail gets a single email.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#single-email
func (s *UsersService) GetEmail(email int, options ...OptionFunc) (*Email, *Response, error) {
	u := fmt.Sprintf("user/emails/%d", email)

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	e := new(Email)
	resp, err := s.client.Do(req, e)
	if err != nil {
		return nil, resp, err
	}

	return e, resp, err
}

// AddEmailOptions represents the available AddEmail() options.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/projects.html#add-email
type AddEmailOptions struct {
	Email *string `url:"email,omitempty" bson:"email,omitempty" json:"email,omitempty"`
}

// AddEmail creates a new email owned by the currently authenticated user.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#add-email
func (s *UsersService) AddEmail(opt *AddEmailOptions, options ...OptionFunc) (*Email, *Response, error) {
	req, err := s.client.NewRequest("POST", "user/emails", opt, options)
	if err != nil {
		return nil, nil, err
	}

	e := new(Email)
	resp, err := s.client.Do(req, e)
	if err != nil {
		return nil, resp, err
	}

	return e, resp, err
}

// AddEmailForUser creates new email owned by specified user. Available only for
// admin.
//
// GitLab API docs: https://docs.gitlab.com/ce/api/users.html#add-email-for-user
func (s *UsersService) AddEmailForUser(user int, opt *AddEmailOptions, options ...OptionFunc) (*Email, *Response, error) {
	u := fmt.Sprintf("users/%d/emails", user)

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	e := new(Email)
	resp, err := s.client.Do(req, e)
	if err != nil {
		return nil, resp, err
	}

	return e, resp, err
}

// DeleteEmail deletes email owned by currently authenticated user. This is an
// idempotent function and calling it on a key that is already deleted or not
// available results in 200 OK.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#delete-email-for-current-owner
func (s *UsersService) DeleteEmail(email int, options ...OptionFunc) (*Response, error) {
	u := fmt.Sprintf("user/emails/%d", email)

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteEmailForUser deletes email owned by a specified user. Available only
// for admin.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#delete-email-for-given-user
func (s *UsersService) DeleteEmailForUser(user, email int, options ...OptionFunc) (*Response, error) {
	u := fmt.Sprintf("users/%d/emails/%d", user, email)

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ImpersonationToken represents an impersonation token.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#get-all-impersonation-tokens-of-a-user
type ImpersonationToken struct {
	ID        int        `bson:"id" json:"id"`
	Name      string     `bson:"name" json:"name"`
	Active    bool       `bson:"active" json:"active"`
	Token     string     `bson:"token" json:"token"`
	Scopes    []string   `bson:"scopes" json:"scopes"`
	Revoked   bool       `bson:"revoked" json:"revoked"`
	CreatedAt *time.Time `bson:"created_at" json:"created_at"`
	ExpiresAt *ISOTime   `bson:"expires_at" json:"expires_at"`
}

// GetAllImpersonationTokensOptions represents the available
// GetAllImpersonationTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#get-all-impersonation-tokens-of-a-user
type GetAllImpersonationTokensOptions struct {
	ListOptions
	State *string `url:"state,omitempty" bson:"state,omitempty" json:"state,omitempty"`
}

// GetAllImpersonationTokens retrieves all impersonation tokens of a user.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#get-all-impersonation-tokens-of-a-user
func (s *UsersService) GetAllImpersonationTokens(user int, opt *GetAllImpersonationTokensOptions, options ...OptionFunc) ([]*ImpersonationToken, *Response, error) {
	u := fmt.Sprintf("users/%d/impersonation_tokens", user)

	req, err := s.client.NewRequest("GET", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var ts []*ImpersonationToken
	resp, err := s.client.Do(req, &ts)
	if err != nil {
		return nil, resp, err
	}

	return ts, resp, err
}

// GetImpersonationToken retrieves an impersonation token of a user.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#get-an-impersonation-token-of-a-user
func (s *UsersService) GetImpersonationToken(user, token int, options ...OptionFunc) (*ImpersonationToken, *Response, error) {
	u := fmt.Sprintf("users/%d/impersonation_tokens/%d", user, token)

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	t := new(ImpersonationToken)
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

// CreateImpersonationTokenOptions represents the available
// CreateImpersonationToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#create-an-impersonation-token
type CreateImpersonationTokenOptions struct {
	Name      *string    `url:"name,omitempty" bson:"name,omitempty" json:"name,omitempty"`
	Scopes    *[]string  `url:"scopes,omitempty" bson:"scopes,omitempty" json:"scopes,omitempty"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" bson:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// CreateImpersonationToken creates an impersonation token.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#create-an-impersonation-token
func (s *UsersService) CreateImpersonationToken(user int, opt *CreateImpersonationTokenOptions, options ...OptionFunc) (*ImpersonationToken, *Response, error) {
	u := fmt.Sprintf("users/%d/impersonation_tokens", user)

	req, err := s.client.NewRequest("POST", u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	t := new(ImpersonationToken)
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

// RevokeImpersonationToken revokes an impersonation token.
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#revoke-an-impersonation-token
func (s *UsersService) RevokeImpersonationToken(user, token int, options ...OptionFunc) (*Response, error) {
	u := fmt.Sprintf("users/%d/impersonation_tokens/%d", user, token)

	req, err := s.client.NewRequest("DELETE", u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UserActivity represents an entry in the user/activities response
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#get-user-activities-admin-only
type UserActivity struct {
	Username       string   `bson:"username" json:"username"`
	LastActivityOn *ISOTime `bson:"last_activity_on" json:"last_activity_on"`
}

// GetUserActivitiesOptions represents the options for GetUserActivities
//
// GitLap API docs:
// https://docs.gitlab.com/ce/api/users.html#get-user-activities-admin-only
type GetUserActivitiesOptions struct {
	From *ISOTime `url:"from,omitempty" bson:"from,omitempty" json:"from,omitempty"`
}

// GetUserActivities retrieves user activities (admin only)
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#get-user-activities-admin-only
func (s *UsersService) GetUserActivities(opt *GetUserActivitiesOptions, options ...OptionFunc) ([]*UserActivity, *Response, error) {
	req, err := s.client.NewRequest("GET", "user/activities", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var t []*UserActivity
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}

	return t, resp, err
}

// UserStatus represents the current status of a user
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#user-status
type UserStatus struct {
	Emoji       string `bson:"emoji" json:"emoji"`
	Message     string `bson:"message" json:"message"`
	MessageHTML string `bson:"message_html" json:"message_html"`
}

// CurrentUserStatus retrieves the user status
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#user-status
func (s *UsersService) CurrentUserStatus(options ...OptionFunc) (*UserStatus, *Response, error) {
	req, err := s.client.NewRequest("GET", "user/status", nil, options)
	if err != nil {
		return nil, nil, err
	}

	status := new(UserStatus)
	resp, err := s.client.Do(req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, err
}

// GetUserStatus retrieves a user's status
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#get-the-status-of-a-user
func (s *UsersService) GetUserStatus(user int, options ...OptionFunc) (*UserStatus, *Response, error) {
	u := fmt.Sprintf("users/%d/status", user)

	req, err := s.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	status := new(UserStatus)
	resp, err := s.client.Do(req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, err
}

// UserStatusOptions represents the options required to set the status
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#set-user-status
type UserStatusOptions struct {
	Emoji   *string `url:"emoji,omitempty" bson:"emoji,omitempty" json:"emoji,omitempty"`
	Message *string `url:"message,omitempty" bson:"message,omitempty" json:"message,omitempty"`
}

// SetUserStatus sets the user's status
//
// GitLab API docs:
// https://docs.gitlab.com/ce/api/users.html#set-user-status
func (s *UsersService) SetUserStatus(opt *UserStatusOptions, options ...OptionFunc) (*UserStatus, *Response, error) {
	req, err := s.client.NewRequest("PUT", "user/status", opt, options)
	if err != nil {
		return nil, nil, err
	}

	status := new(UserStatus)
	resp, err := s.client.Do(req, status)
	if err != nil {
		return nil, resp, err
	}

	return status, resp, err
}
