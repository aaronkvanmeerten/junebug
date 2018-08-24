package junebug

const (
	prompt = `{
    "text": "JuneBug here to help you setup your meeting!",
    "attachments": [
        {
            "text": "Need video?",
            "fallback": "You will not have any remote attendance",
            "callback_id": "%s",
            "color": "#3AA3E3",
            "attachment_type": "default",
            "actions": [
                {
                    "name": "video",
                    "text": "Jitsi Meet",
                    "type": "button",
                    "value": "jitsi"
                },
                {
                    "name": "video",
                    "text": "Bluejeans",
                    "type": "button",
                    "value": "bluejeans"
                }
            ]
        },
        {
            "text": "Need to track minutes?",
            "fallback": "You will not be able to track minutes",
            "callback_id": "%s",
            "color": "#34E2E0",
            "attachment_type": "default",
            "actions": [
                {
                    "name": "tracking",
                    "text": "Confluence",
                    "type": "button",
                    "value": "confluence"
                },
                {
                    "name": "tracking",
                    "text": "Trello",
                    "type": "button",
                    "value": "trello"
                }
            ]
        },
        {
            "text": "Whenever you're ready to share with the channel...",
            "fallback": "You will not meeting today.",
            "callback_id": "%s",
            "color": "#6BD838",
            "attachment_type": "default",
            "actions": [
                {
                    "name": "start",
                    "text": "Start",
                    "type": "button",
                    "value": "start"
                }
            ]
        }
    ]
}`
	roomTemplate = `{"response_type":"in_channel","attachments":[{"fallback":"Meeting started %s","title":"Meeting started %s","color":"#3AA3E3","attachment_type":"default","actions":[{"name":"join","text":"Join","type":"button","url":"%s","style":"primary"}]}]}`
)
