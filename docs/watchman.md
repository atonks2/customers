[Moov Watchman](https://github.com/moov-io/watchman) is used by Customers to compare Customer information against Anit-Money Laundering (AML) and sanction lists provided by various government agencies. The primary list checked is OFAC as a component of Know Your Customer (KYC). These checks are required of all US businesses.

### Configuration

Watchman offers a few [environmental variables](https://github.com/moov-io/watchman#configuration) for reading the lists, search tuning, and binding to different HTTP ports.

### OFAC Checks

As required by United States law and NACHA guidelines all transfers are checked against the Office of Foreign Asset Control (OFAC) lists for sanctioned individuals and entities to combat fraud, terrorism and unlawful monetary transfers outside of the United States. Customers uses Watchman to perform these checks.

### OFAC Searches

Customers performs searches against the OFAC list of entities which US businesses are blocked from doing business with. This list changes frequently with world politics and policies. Customer objects are required to be in either `ReceiveOnly` or `Validated` status in order for `Transfers` to be created.
