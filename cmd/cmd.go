// Copyright (c) 2020 Mercari, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cloudspannerecosystem/wrench/pkg/spanner"
)

const (
	flagNameProject               = "project"
	flagNameInstance              = "instance"
	flagNameDatabase              = "database"
	flagNameDirectory             = "directory"
	flagCredentialsFile           = "credentials_file"
	flagNameMigrationTable        = "migration_table"
	flagProtoDescriptorsFile      = "proto_descriptors_file"
	flagImpersonateServiceAccount = "impersonate_service_account"
	flagNameSchemaFile            = "schema_file"
	flagDDLFile                   = "ddl"
	flagDMLFile                   = "dml"
	flagPartitioned               = "partitioned"
	flagPriority                  = "priority"
	flagNode                      = "node"
	flagTimeout                   = "timeout"
	defaultSchemaFileName         = "schema.sql"
	defaultMigrationTableName     = "SchemaMigrations"
)

func newSpannerClient(ctx context.Context, c *cobra.Command) (*spanner.Client, error) {
	config := &spanner.Config{
		Project:                   c.Flag(flagNameProject).Value.String(),
		Instance:                  c.Flag(flagNameInstance).Value.String(),
		Database:                  c.Flag(flagNameDatabase).Value.String(),
		CredentialsFile:           c.Flag(flagCredentialsFile).Value.String(),
		ImpersonateServiceAccount: c.Flag(flagImpersonateServiceAccount).Value.String(),
	}

	desc, err := protoDescriptors(c)
	if err != nil {
		return nil, err
	}
	config.ProtoDescriptors = desc

	client, err := spanner.NewClient(ctx, config)
	if err != nil {
		return nil, &Error{
			err: err,
			cmd: c,
		}
	}

	return client, nil
}

func newSpannerAdminClient(ctx context.Context, c *cobra.Command) (*spanner.AdminClient, error) {
	config := &spanner.Config{
		Project:                   c.Flag(flagNameProject).Value.String(),
		Instance:                  c.Flag(flagNameInstance).Value.String(),
		Database:                  c.Flag(flagNameDatabase).Value.String(),
		CredentialsFile:           c.Flag(flagCredentialsFile).Value.String(),
		ImpersonateServiceAccount: c.Flag(flagImpersonateServiceAccount).Value.String(),
	}

	desc, err := protoDescriptors(c)
	if err != nil {
		return nil, err
	}
	config.ProtoDescriptors = desc

	client, err := spanner.NewAdminClient(ctx, config)
	if err != nil {
		return nil, &Error{
			err: err,
			cmd: c,
		}
	}

	return client, nil
}

func schemaFilePath(c *cobra.Command) string {
	filename := c.Flag(flagNameSchemaFile).Value.String()
	if filename == "" {
		filename = defaultSchemaFileName
	}
	return filepath.Join(c.Flag(flagNameDirectory).Value.String(), filename)
}

func protoDescriptors(c *cobra.Command) ([]byte, error) {
	filename := c.Flag(flagProtoDescriptorsFile).Value.String()
	if filename == "" {
		return nil, nil
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, &Error{
			err: err,
			cmd: c,
		}
	}
	return b, nil
}
