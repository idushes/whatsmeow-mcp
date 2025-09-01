# WhatsApp MCP Tools Implementation Plan

This document describes all planned tools for the whatsmeow-mcp server and tracks their implementation status. Each tool corresponds to WhatsApp functionality provided by the whatsmeow Go library.

## Implementation Status Legend
- ✅ **Implemented** - Tool is fully implemented and tested
- 🚧 **In Progress** - Tool is currently being implemented
- ⏳ **Planned** - Tool is planned for implementation
- ❌ **Blocked** - Tool implementation is blocked by dependencies

## Implementation Progress Summary
**Total Tools:** 38  
**Implemented:** 7 (18%)  
**In Progress:** 0 (0%)  
**Planned:** 31 (82%)  
**Blocked:** 0 (0%)

## Quick Tool Index

### Connection and Authentication Tools (3 tools)
- [`get_qr_code`](#get_qr_code-) ✅ - Generate QR code for WhatsApp Web authentication
- [`logout`](#logout-) ⏳ - Logout from WhatsApp account
- [`is_logged_in`](#is_logged_in-) ✅ - Check WhatsApp authentication status

### Message Sending Tools (10 tools)
- [`send_message`](#send_message-) ✅ - Send a text message to a WhatsApp chat or contact
- [`send_image_message`](#send_image_message-) ⏳ - Send image with optional caption
- [`send_document_message`](#send_document_message-) ⏳ - Send document/file
- [`send_audio_message`](#send_audio_message-) ⏳ - Send audio message
- [`send_video_message`](#send_video_message-) ⏳ - Send video message
- [`send_location_message`](#send_location_message-) ⏳ - Send location message
- [`build_poll_creation`](#build_poll_creation-) ⏳ - Create a poll message
- [`build_poll_vote`](#build_poll_vote-) ⏳ - Vote in a poll
- [`build_reaction`](#build_reaction-) ⏳ - Add reaction to a message
- [`build_edit`](#build_edit-) ⏳ - Edit a previously sent message
- [`build_revoke`](#build_revoke-) ⏳ - Revoke/delete a sent message

### Group Management Tools (8 tools)
- [`create_group`](#create_group-) ⏳ - Create new WhatsApp group
- [`get_group_info`](#get_group_info-) ⏳ - Get detailed group information
- [`join_group_with_link`](#join_group_with_link-) ⏳ - Join group using invite link
- [`join_group_with_invite`](#join_group_with_invite-) ⏳ - Join group using invite message
- [`leave_group`](#leave_group-) ⏳ - Leave a group
- [`set_group_name`](#set_group_name-) ⏳ - Change group name
- [`set_group_description`](#set_group_description-) ⏳ - Change group description
- [`set_group_photo`](#set_group_photo-) ⏳ - Set group profile photo
- [`update_group_participants`](#update_group_participants-) ⏳ - Add or remove group participants

### Contact and User Information Tools (6 tools)
- [`get_user_info`](#get_user_info-) ⏳ - Get user information including avatar, status, and verification
- [`get_user_devices`](#get_user_devices-) ⏳ - Get list of user's devices
- [`is_on_whatsapp`](#is_on_whatsapp-) ✅ - Check if phone numbers are registered on WhatsApp and get their JIDs
- [`get_profile_picture_info`](#get_profile_picture_info-) ⏳ - Get profile picture information
- [`get_business_profile`](#get_business_profile-) ⏳ - Get business profile information
- [`get_contacts`](#get_contacts-) ⏳ - Get list of contacts

### Profile Management Tools (1 tool)
- [`set_status_message`](#set_status_message-) ⏳ - Set user status message

### Presence and Status Tools (3 tools)
- [`send_presence`](#send_presence-) ⏳ - Set global presence status
- [`subscribe_presence`](#subscribe_presence-) ⏳ - Subscribe to user's presence updates
- [`send_chat_presence`](#send_chat_presence-) ⏳ - Send typing or recording status to specific chat

### Chat Management Tools (1 tool)
- [`get_all_chats`](#get_all_chats-) ⏳ - Get list of all chats

### Media Tools (2 tools)
- [`download`](#download-) ⏳ - Download media from message
- [`upload`](#upload-) ⏳ - Upload media file to WhatsApp servers

### Message Management Tools (3 tools)
- [`mark_read`](#mark_read-) ⏳ - Mark messages as read
- [`mark_messages_as_read`](#mark_messages_as_read-) ✅ - Mark all unread messages in a chat as read
- [`get_chat_history`](#get_chat_history-) ✅ - Retrieve message history from a WhatsApp conversation with pagination support

### Notification Tools (1 tool)
- [`get_unread_messages`](#get_unread_messages-) ✅ - Retrieve unread messages from WhatsApp chats

### Privacy and Settings Tools (2 tools)
- [`get_privacy_settings`](#get_privacy_settings-) ⏳ - Get current privacy settings
- [`get_blocklist`](#get_blocklist-) ⏳ - Get list of blocked contacts



---

## Detailed Tool Descriptions

## Connection and Authentication Tools


### `get_qr_code` ✅
**Status:** Implemented  
**Description:** Generate QR code for WhatsApp Web authentication  
**Parameters:**
- None

**Returns:**
- `qr_code`: string - Raw QR code string content
- `code`: string - Same as qr_code (for compatibility)
- `image_url`: string - URL to hosted QR code image file
- `timeout`: number - QR code expiration time in seconds (30)
- `expires_at`: number - Unix timestamp when QR code expires
- `success`: boolean - Operation success status



### `logout` ⏳
**Status:** Planned  
**Description:** Logout from WhatsApp account  
**Parameters:**
- None

**Returns:**
- `success`: boolean - Logout status
- `message`: string - Status message



### `is_logged_in` ✅
**Status:** Implemented  
**Description:** Check WhatsApp authentication status  
**Parameters:**
- None

**Returns:**
- `logged_in`: boolean - Authentication status
- `success`: boolean - Request status

## Message Sending Tools

### `send_message` ✅
**Status:** Implemented  
**Description:** Send a text message to a WhatsApp chat or contact. Requires authentication.  
**Parameters:**
- `to`: string - Recipient JID (e.g., "1234567890@s.whatsapp.net" for contact, "1234567890-1234567890@g.us" for group)
- `text`: string - Message text content
- `quoted_message_id`: string (optional) - ID of message to quote/reply to

**Returns:**
- `message_id`: string - Sent message ID
- `timestamp`: number - Message timestamp (Unix timestamp)
- `success`: boolean - Send status
- `to`: string - Recipient JID (echoed back)
- `text`: string - Message text (echoed back)
- `quoted_message_id`: string (optional) - Quoted message ID if provided

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

### `send_location_message` ⏳
**Status:** Planned  
**Description:** Send location message  
**Parameters:**
- `to`: string - Recipient JID
- `latitude`: number - Location latitude
- `longitude`: number - Location longitude
- `name`: string (optional) - Location name/title
- `address`: string (optional) - Location address

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

### `is_on_whatsapp` ✅
**Status:** Implemented  
**Description:** Check if phone numbers are registered on WhatsApp and get their JIDs.  
**Parameters:**
- `phones`: array of strings - Phone numbers in international format

**Returns:**
- `results`: array of objects - Registration status for each phone number
  - `phone`: string - Original phone number
  - `is_on_whatsapp`: boolean - Registration status
  - `jid`: string - WhatsApp JID if registered
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

### `get_contacts` ⏳
**Status:** Planned  
**Description:** Get list of contacts  
**Parameters:**
- None

**Returns:**
- `contacts`: array - Array of contact objects with JID, name, and status
- `success`: boolean - Request status

## Profile Management Tools

### `set_status_message` ⏳
**Status:** Planned  
**Description:** Set user status message  
**Parameters:**
- `status`: string - Status message text

**Returns:**
- `success`: boolean - Update status

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

### `mark_messages_as_read` ✅
**Status:** Implemented  
**Description:** Mark all unread messages in a specific chat as read. This tool updates the read status of all unread messages in the specified chat.  
**Parameters:**
- `chat`: string - WhatsApp JID (chat identifier) to mark messages as read in this chat. For phone numbers: 'phonenumber@s.whatsapp.net' (e.g. '1234567890@s.whatsapp.net'). For groups: 'groupid@g.us'

**Returns:**
- `success`: boolean - Operation success status
- `chat`: string - Chat JID (echoed back)
- `message`: string - Descriptive message about the operation result

### `get_chat_history` ✅
**Status:** Implemented  
**Description:** Retrieve message history from a WhatsApp conversation with pagination support.  
**Parameters:**
- `chat`: string - Chat JID
- `count`: number (optional) - Number of messages to retrieve (default: 50, max: 100)
- `before_message_id`: string (optional) - Get messages before this ID

**Returns:**
- `messages`: array of objects - Array of message objects
  - `id`: string - Message ID
  - `from`: string - Sender JID
  - `to`: string - Recipient JID (optional)
  - `text`: string - Message text content
  - `timestamp`: number - Unix timestamp
  - `chat`: string - Chat JID
  - `quoted_message_id`: string (optional) - ID of quoted message
- `has_more`: boolean - Whether more messages are available
- `success`: boolean - Request status
- `chat`: string - Chat JID (echoed back)
- `count`: number - Actual number of messages returned

## Chat Management Tools

### `get_all_chats` ⏳
**Status:** Planned  
**Description:** Get list of all chats  
**Parameters:**
- None

**Returns:**
- `chats`: array - Array of chat objects with JID, name, type, and last message info
- `success`: boolean - Request status

## Notification Tools

### `get_unread_messages` ✅
**Status:** Implemented  
**Description:** Retrieve unread messages from WhatsApp chats.  
**Parameters:**
- `chat`: string (optional) - WhatsApp JID to filter unread messages from a specific chat. If omitted, returns unread messages from all chats
- `count`: number (optional) - Maximum number of unread messages to retrieve (default: 50, max: 100)

**Returns:**
- `messages`: array of objects - Array of unread message objects
  - `id`: string - Message ID
  - `from`: string - Sender JID
  - `to`: string - Recipient JID (optional)
  - `text`: string - Message text content
  - `timestamp`: number - Unix timestamp
  - `chat`: string - Chat JID
  - `quoted_message_id`: string (optional) - ID of quoted message
- `success`: boolean - Request status
- `chat`: string (optional) - Chat JID filter (echoed back if provided)
- `count`: number - Actual number of messages returned

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
