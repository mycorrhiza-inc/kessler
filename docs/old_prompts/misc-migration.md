Scraping the Ny PUC depends on 
- Getting Marker Working.
- Database Migration Working.


View for individual dockets


Database 
- Files with the following fields 
Name
Date Published 
File Extension
isPrivate
hash 
lang

- Document Text Source (UNCHANGED)

- File Metadata Table With the following fields:
Object ID (foreign key to files)
Metadata Field 
Metadata Value (Jsonified, or compressed in some way to store things like lists, or dates, or numbers)

- Source Table with the following fields
Object ID (foreign key to files)
Source Type (type of source, url, doi, )
URL (Optional)
...Other stuff as needed to decide later

- Juristiction Table (For filtering on juristiction, delay until later.)


- Stage Table with the following fields
Object ID (foreign key to files)
stage (string enum with same values as in python)
error (optional string)




- Other tables with the following points behind them

- 




Misc Random:
Hash Deduplication - Implement if we ever upload duplicate files.


