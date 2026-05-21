# DynamoDB Single-Table Design

## Table: LinkTable
- **PK**: USER#<userID>
- **SK**: BOOKMARK#<bookmarkID>

## Access Patterns
1. Get all bookmarks for user → Query PK=USER#X, SK begins_with BOOKMARK#
2. Get specific bookmark → GetItem PK+SK
3. Delete bookmark → DeleteItem PK+SK

## Future: Tags GSI
- GSI1PK: TAG#<tagName>
- GSI1SK: BOOKMARK#<bookmarkID>
- Query by tag → Query GSI1