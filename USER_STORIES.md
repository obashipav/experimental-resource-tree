# User Stories

## Deleting things from the path.

Deleting data from a service is a big topic by its own. We want to use `user stories` and try to see what will happen if we try to delete things 
that have dependencies. 

## User stories
 We will group the user stories per resource

### Folder stories
As a user | I want to | How to test | Notes
--- | --- | ---- | ---- 
Who has admin privileges | delete a folder that is empty | Create a  folder that is empty and then delete it | 
Who has admin privileges | delete a folder that has content | Create a folder, then create other resources under the folder and try to delete it | 

Let's assume the following structure

Folder: `Dataflow Unit 2021`
- Content: `empty` 

_with content_

Assume the following structure, (and we don't take into account assets, templates, bit, dav, fields)

- Org: `Forth Valley College`
    - Folder: `Dataflow Unit 2021`
        - Folder: `Solutions`
        - Folder: `Markings`
            - Folder: `Alice GID 1392`
                - Project: `Spring Assignment 2021 GID 1392`
            - Folder: `Bob GID 2201`
                - Project: `Spring Assignment 2021 GID 1392`
        - Project: `Dataflow Forth Valley`
    - Project: `Public Unit 2021`
    - Repository: `Forth Valley Assets`
    - Repository: `Dataflow Unit Assets`
    
__How do we delete a resource that has content inside?__
    