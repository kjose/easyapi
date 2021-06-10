// Copyright 2021 Kévin José.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

const (
	// Resource events
	EVENT_RESOURCE_POST_READ   = "resource.post_read"
	EVENT_RESOURCE_PRE_CREATE  = "resource.pre_create"
	EVENT_RESOURCE_POST_CREATE = "resource.post_create"
	EVENT_RESOURCE_PRE_UPDATE  = "resource.pre_update"
	EVENT_RESOURCE_POST_UPDATE = "resource.post_update"
	EVENT_RESOURCE_PRE_DELETE  = "resource.pre_delete"
	EVENT_RESOURCE_ACTION      = "resource.action"

	// Request events
	EVENT_REQUEST_START     = "request.start"
	EVENT_REQUEST_TERMINATE = "request.terminate"
)
