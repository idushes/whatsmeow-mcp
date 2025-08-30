# WhatsApp MCP Tools Implementation Plan

This document describes all planned tools for the whatsmeow-mcp server and tracks their implementation status. Each tool corresponds to WhatsApp functionality provided by the whatsmeow Go library.

## Implementation Status Legend
- ✅ **Implemented** - Tool is fully implemented and tested
- 🚧 **In Progress** - Tool is currently being implemented
- ⏳ **Planned** - Tool is planned for implementation
- ❌ **Blocked** - Tool implementation is blocked by dependencies

## Implementation Progress Summary
**Total Tools:** 38  
**Implemented:** 0 (0%)  
**In Progress:** 0 (0%)  
**Planned:** 38 (100%)  
**Blocked:** 0 (0%)

## Connection and Authentication Tools

### `connect` ⏳
**Status:** Planned  
**Description:** Establishes connection to WhatsApp servers  
**Parameters:**
- None

**Returns:**
- `success`: boolean - Connection status
- `message`: string - Status message

**Example:**
```json
{
  "name": "connect",
  "arguments": {}
}
```

### `disconnect` ⏳
**Status:** Planned  
**Description:** Disconnects from WhatsApp servers  
**Parameters:**
- None

**Returns:**
- `success`: boolean - Disconnection status
- `message`: string - Status message

### `get_qr_code` ⏳
**Status:** Planned  
**Description:** Generates QR code for WhatsApp Web login  
**Parameters:**
- None

**Returns:**
- `qr_code`: string - Base64 encoded QR code image
- `code`: string - QR code text content
- `timeout`: number - QR code expiration time in seconds

### `pair_phone` ⏳
**Status:** Planned  
**Description:** Pair device using phone number  
**Parameters:**
- `phone`: string - Phone number in international format (e.g., "+1234567890")
- `show_push_notification`: boolean - Whether to show push notification

**Returns:**
- `pairing_code`: string - 8-digit pairing code
- `success`: boolean - Pairing initiation status

### `logout` ⏳
**Status:** Planned  
**Description:** Logout from WhatsApp account  
**Parameters:**
- None

**Returns:**
- `success`: boolean - Logout status
- `message`: string - Status message

### `is_connected` ⏳
**Status:** Planned  
**Description:** Check if client is connected to WhatsApp servers  
**Parameters:**
- None

**Returns:**
- `connected`: boolean - Connection status

### `is_logged_in` ⏳
**Status:** Planned  
**Description:** Check if user is authenticated  
**Parameters:**
- None

**Returns:**
- `logged_in`: boolean - Authentication status

## Message Sending Tools

### `send_message` ⏳
**Status:** Planned  
**Description:** Send text message to chat or contact  
**Parameters:**
- `to`: string - Recipient JID (e.g., "1234567890@s.whatsapp.net" for contact, "1234567890-1234567890@g.us" for group)
- `text`: string - Message text content
- `quoted_message_id`: string (optional) - ID of message to quote/reply to

**Returns:**
- `message_id`: string - Sent message ID
- `timestamp`: number - Message timestamp
- `success`: boolean - Send status

### `send_image_message` ⏳
**Status:** Planned  
**Description:** Send image with optional caption  
**Parameters:**
- `to`: string - Recipient JID
- `image_path`: string - Path to image file
- `caption`: string (optional) - Image caption
- `quoted_message_id`: string (optional) - ID of message to quote/reply to

**Returns:**
- `message_id`: string - Sent message ID
- `timestamp`: number - Message timestamp
- `success`: boolean - Send status

### `send_document_message` ⏳
**Status:** Planned  
**Description:** Send document/file  
**Parameters:**
- `to`: string - Recipient JID
- `document_path`: string - Path to document file
- `filename`: string (optional) - Custom filename
- `mimetype`: string (optional) - Document MIME type
- `caption`: string (optional) - Document caption

**Returns:**
- `message_id`: string - Sent message ID
- `timestamp`: number - Message timestamp
- `success`: boolean - Send status

### `send_audio_message` ⏳
**Status:** Planned  
**Description:** Send audio message  
**Parameters:**
- `to`: string - Recipient JID
- `audio_path`: string - Path to audio file
- `ptt`: boolean (optional) - Whether audio is push-to-talk/voice note

**Returns:**
- `message_id`: string - Sent message ID
- `timestamp`: number - Message timestamp
- `success`: boolean - Send status

### `send_video_message` ⏳
**Status:** Planned  
**Description:** Send video message  
**Parameters:**
- `to`: string - Recipient JID
- `video_path`: string - Path to video file
- `caption`: string (optional) - Video caption

**Returns:**
- `message_id`: string - Sent message ID
- `timestamp`: number - Message timestamp
- `success`: boolean - Send status

### `build_poll_creation` ⏳
**Status:** Planned  
**Description:** Create a poll message  
**Parameters:**
- `to`: string - Recipient JID
- `name`: string - Poll question
- `options`: array of strings - Poll options
- `selectable_count`: number (optional) - Number of options user can select (default: 1)

**Returns:**
- `message_id`: string - Sent message ID
- `timestamp`: number - Message timestamp
- `success`: boolean - Send status

### `build_poll_vote` ⏳
**Status:** Planned  
**Description:** Vote in a poll  
**Parameters:**
- `chat`: string - Chat JID where poll is located
- `poll_message_id`: string - Poll message ID
- `option_names`: array of strings - Selected option names

**Returns:**
- `message_id`: string - Vote message ID
- `success`: boolean - Vote status

### `build_reaction` ⏳
**Status:** Planned  
**Description:** Add reaction to a message  
**Parameters:**
- `chat`: string - Chat JID
- `message_id`: string - Target message ID
- `sender`: string - Original message sender JID
- `reaction`: string - Reaction emoji (empty string to remove reaction)

**Returns:**
- `message_id`: string - Reaction message ID
- `success`: boolean - Reaction status

### `build_edit` ⏳
**Status:** Planned  
**Description:** Edit a previously sent message  
**Parameters:**
- `chat`: string - Chat JID
- `message_id`: string - Message ID to edit
- `new_text`: string - New message content

**Returns:**
- `message_id`: string - Edit message ID
- `success`: boolean - Edit status

### `build_revoke` ⏳
**Status:** Planned  
**Description:** Revoke/delete a sent message  
**Parameters:**
- `chat`: string - Chat JID
- `message_id`: string - Message ID to revoke
- `sender`: string - Original message sender JID

**Returns:**
- `message_id`: string - Revoke message ID
- `success`: boolean - Revoke status

## Group Management Tools

### `create_group` ⏳
**Status:** Planned  
**Description:** Create new WhatsApp group  
**Parameters:**
- `name`: string - Group name
- `participants`: array of strings - Participant JIDs
- `description`: string (optional) - Group description

**Returns:**
- `group_jid`: string - Created group JID
- `success`: boolean - Creation status
- `group_info`: object - Group information

### `get_group_info` ⏳
**Status:** Planned  
**Description:** Get detailed group information  
**Parameters:**
- `group_jid`: string - Group JID

**Returns:**
- `group_info`: object - Group details including participants, admins, settings
- `success`: boolean - Request status

### `join_group_with_link` ⏳
**Status:** Planned  
**Description:** Join group using invite link  
**Parameters:**
- `invite_code`: string - Group invite code from link

**Returns:**
- `group_jid`: string - Joined group JID
- `success`: boolean - Join status

### `join_group_with_invite` ⏳
**Status:** Planned  
**Description:** Join group using invite message  
**Parameters:**
- `group_jid`: string - Group JID
- `inviter`: string - Inviter JID
- `code`: string - Invite code
- `expiration`: number - Invite expiration timestamp

**Returns:**
- `success`: boolean - Join status

### `leave_group` ⏳
**Status:** Planned  
**Description:** Leave a group  
**Parameters:**
- `group_jid`: string - Group JID to leave

**Returns:**
- `success`: boolean - Leave status

### `set_group_name` ⏳
**Status:** Planned  
**Description:** Change group name  
**Parameters:**
- `group_jid`: string - Group JID
- `name`: string - New group name

**Returns:**
- `success`: boolean - Update status

### `set_group_description` ⏳
**Status:** Planned  
**Description:** Change group description  
**Parameters:**
- `group_jid`: string - Group JID
- `description`: string - New group description

**Returns:**
- `success`: boolean - Update status

### `set_group_photo` ⏳
**Status:** Planned  
**Description:** Set group profile photo  
**Parameters:**
- `group_jid`: string - Group JID
- `image_path`: string - Path to image file

**Returns:**
- `picture_id`: string - New picture ID
- `success`: boolean - Update status

### `update_group_participants` ⏳
**Status:** Planned  
**Description:** Add or remove group participants  
**Parameters:**
- `group_jid`: string - Group JID
- `participants`: array of strings - Participant JIDs
- `action`: string - "add", "remove", "promote", or "demote"

**Returns:**
- `results`: array - Results for each participant
- `success`: boolean - Overall operation status

## Contact and User Information Tools

### `get_user_info` ⏳
**Status:** Planned  
**Description:** Get user information including avatar, status, and verification  
**Parameters:**
- `jids`: array of strings - User JIDs to query

**Returns:**
- `users`: object - Map of JID to user information
- `success`: boolean - Request status

### `get_user_devices` ⏳
**Status:** Planned  
**Description:** Get list of user's devices  
**Parameters:**
- `jids`: array of strings - User JIDs to query

**Returns:**
- `devices`: array of strings - Device JIDs
- `success`: boolean - Request status

### `is_on_whatsapp` ⏳
**Status:** Planned  
**Description:** Check if phone numbers are registered on WhatsApp  
**Parameters:**
- `phones`: array of strings - Phone numbers in international format

**Returns:**
- `results`: array - Registration status for each phone number
- `success`: boolean - Request status

### `get_profile_picture_info` ⏳
**Status:** Planned  
**Description:** Get profile picture information  
**Parameters:**
- `jid`: string - User JID
- `preview`: boolean (optional) - Whether to get preview quality

**Returns:**
- `picture_info`: object - Picture URL, ID, and metadata
- `success`: boolean - Request status

### `get_business_profile` ⏳
**Status:** Planned  
**Description:** Get business profile information  
**Parameters:**
- `jid`: string - Business JID

**Returns:**
- `business_profile`: object - Business information
- `success`: boolean - Request status

## Presence and Status Tools

### `send_presence` ⏳
**Status:** Planned  
**Description:** Set global presence status  
**Parameters:**
- `presence`: string - Presence state: "available", "unavailable"

**Returns:**
- `success`: boolean - Update status

### `subscribe_presence` ⏳
**Status:** Planned  
**Description:** Subscribe to user's presence updates  
**Parameters:**
- `jid`: string - User JID to subscribe to

**Returns:**
- `success`: boolean - Subscription status

### `send_chat_presence` ⏳
**Status:** Planned  
**Description:** Send typing or recording status to specific chat  
**Parameters:**
- `jid`: string - Chat JID
- `state`: string - Presence state: "composing", "recording", "paused"
- `media`: string (optional) - Media type for recording state

**Returns:**
- `success`: boolean - Send status

## Media Tools

### `download` ⏳
**Status:** Planned  
**Description:** Download media from message  
**Parameters:**
- `message`: object - Message object containing media info
- `save_path`: string (optional) - Path to save downloaded file

**Returns:**
- `file_path`: string - Path to downloaded file
- `file_size`: number - File size in bytes
- `success`: boolean - Download status

### `upload` ⏳
**Status:** Planned  
**Description:** Upload media file to WhatsApp servers  
**Parameters:**
- `file_path`: string - Path to file to upload
- `media_type`: string - Media type: "image", "video", "audio", "document"

**Returns:**
- `media_key`: string - Media encryption key
- `file_sha256`: string - File hash
- `file_enc_sha256`: string - Encrypted file hash
- `direct_path`: string - Media direct path
- `success`: boolean - Upload status

## Message Management Tools

### `mark_read` ⏳
**Status:** Planned  
**Description:** Mark messages as read  
**Parameters:**
- `chat`: string - Chat JID
- `message_ids`: array of strings - Message IDs to mark as read
- `sender`: string (optional) - Message sender JID

**Returns:**
- `success`: boolean - Mark read status

### `get_chat_history` ⏳
**Status:** Planned  
**Description:** Get chat message history  
**Parameters:**
- `chat`: string - Chat JID
- `count`: number (optional) - Number of messages to retrieve (default: 50)
- `before_message_id`: string (optional) - Get messages before this ID

**Returns:**
- `messages`: array - Array of message objects
- `has_more`: boolean - Whether more messages are available
- `success`: boolean - Request status

## Privacy and Settings Tools

### `get_privacy_settings` ⏳
**Status:** Planned  
**Description:** Get current privacy settings  
**Parameters:**
- None

**Returns:**
- `privacy_settings`: object - Privacy configuration
- `success`: boolean - Request status

### `get_blocklist` ⏳
**Status:** Planned  
**Description:** Get list of blocked contacts  
**Parameters:**
- None

**Returns:**
- `blocked_contacts`: array of strings - Blocked JIDs
- `success`: boolean - Request status

## Utility Tools

### `generate_message_id` ⏳
**Status:** Planned  
**Description:** Generate unique message ID  
**Parameters:**
- None

**Returns:**
- `message_id`: string - Generated message ID

### `parse_jid` ⏳
**Status:** Planned  
**Description:** Parse and validate JID format  
**Parameters:**
- `jid`: string - JID to parse

**Returns:**
- `user`: string - User part of JID
- `server`: string - Server part of JID
- `device`: number - Device ID (for AD JIDs)
- `is_group`: boolean - Whether JID is a group
- `is_broadcast`: boolean - Whether JID is a broadcast list
- `valid`: boolean - Whether JID format is valid

## Error Handling

All tools return a standardized error format when operations fail:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Additional error details if available"
  }
}
```

## Common Error Codes

- `NOT_CONNECTED`: Client is not connected to WhatsApp
- `NOT_LOGGED_IN`: Client is not authenticated
- `INVALID_JID`: Invalid JID format
- `RATE_LIMITED`: Too many requests
- `MEDIA_UPLOAD_FAILED`: Media upload failed
- `MESSAGE_SEND_FAILED`: Message sending failed
- `GROUP_NOT_FOUND`: Group does not exist
- `INSUFFICIENT_PERMISSIONS`: User lacks required permissions
- `NETWORK_ERROR`: Network connectivity issue
