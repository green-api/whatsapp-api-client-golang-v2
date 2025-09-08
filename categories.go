package greenapi

// Account category presents methods for working with the account.
//
// https://green-api.com/v3/docs/api/account/
func (c *GreenAPI) Account() AccountCategory {
	return AccountCategory{GreenAPI: c}
}

// Sending category presents methods for sending different messages.
//
// https://green-api.com/v3/docs/api/sending/
func (c *GreenAPI) Sending() SendingCategory {
	return SendingCategory{GreenAPI: c}
}

// Receiving category presents methods for working with receiving events.
//
// https://green-api.com/v3/docs/api/receiving/
func (c *GreenAPI) Receiving() ReceivingCategory {
	return ReceivingCategory{GreenAPI: c}
}

// Groups category presents methods for working with group chats.
//
// https://green-api.com/v3/docs/api/groups/
func (c *GreenAPI) Groups() GroupsCategory {
	return GroupsCategory{GreenAPI: c}
}

// Journals present methods for working with incoming and outgoing messages.
//
// https://green-api.com/v3/docs/api/journals/
func (c *GreenAPI) Journals() JournalsCategory {
	return JournalsCategory{GreenAPI: c}
}

// Queues category presents methods for working with a messages queue.
//
// https://green-api.com/v3/docs/api/queues/
func (c *GreenAPI) Queues() QueuesCategory {
	return QueuesCategory{GreenAPI: c}
}

// ReadMark category presents methods for working with chat read mark.
//
// https://green-api.com/v3/docs/api/marks/
func (c *GreenAPI) ReadMark() ReadMarkCategory {
	return ReadMarkCategory{GreenAPI: c}
}

// Service category presents different service methods.
//
// https://green-api.com/v3/docs/api/service/
func (c *GreenAPI) Service() ServiceCategory {
	return ServiceCategory{GreenAPI: c}
}

// Partner category presents exclusive methods for partners.
// The partnership scheme involves deeper integration with the service
// and working with a larger number of instances on your side:
//
// * Instance management via API
//
// * Postpaid billing system (starting from the second month of operation)
//
// * Daily billing (for created and not deleted instances)
//
// * Dedicated support line
//
// For questions regarding connection to the partnership scheme
// and additional conditions, please contact us via email
// at support@green-api.com or via chat on the website.
//
// https://green-api.com/v3/docs/partners/
func (c *GreenAPIPartner) Partner() PartnerCategory {
	return PartnerCategory{GreenAPIPartner: c}
}
