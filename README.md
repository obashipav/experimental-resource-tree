# V1 Platform path prototype

The prototype intends to improve the usability of the old matpath and of course improve how we keep resources to the V1
path.

## API endpoint design

Like every prototype we will try to list every endpoint which is going to be used and also will add value to our API.

HTTP Method | URI | Description
--- | --- | ---
GET | /activity | Will return the latest user activity paginated, with default sort by time.
POST | /org | Will create a new resource of type `org` and will return the full URL
GET | /org/{resourceId} | Returns the resource data and other links as full URLs
UPDATE | /org/{resourceId} | Will update the resource data
DELETE | /org/{resourceId} | Will close the organisation and EVERYTHING within, (*snapped-fingers sound in the background)
GET | /orgs | Get a list of all the connected organisations sorted by latest visited as default and paginated, of course it will support query parameters for search





## Resource path is the old `matpath`

This document will focus on the technical side, and we will start from the postgres model.

### Postgres model

Let's start with the resource path model.

```postgresql
--  Table Definitions
CREATE TABLE base
(
    id           uuid    not null primary key default gen_random_uuid(),
    path_uri     varchar not null,
    previous_uri varchar,
    type         varchar not null,
    hierarchy    varchar[],
    hlevel       int,
    label        varchar,
    alt_label    varchar,
    description  varchar,
    tsearch      tsvector generated always as ( label_and_desc_tsearch(label, description) ) stored,
    owner        varchar not null,
    updated_by   varchar,
    created_at   bigint default nownano(CURRENT_TIMESTAMP),
    updated_at   bigint default nownano(CURRENT_TIMESTAMP),
    deleted      bool default false,
    constraint pathing unique (path_uri, label, type)
);

CREATE TABLE org
(
    LIKE base including all
) inherits (base);

CREATE TABLE repository
(
    LIKE base including all
) inherits (base);

CREATE TABLE folder
(
    LIKE base including all
) INHERITS (base);

CREATE TABLE project
(
    LIKE base including all,
    colour varchar
) INHERITS (base);
```

#### Primary key Path URL

First, the primary key if the resource's path URL. The URL is generated using the `shortuuid` library with a `namespace`
. The table is looking like this:

path_uri | resource_id | type | previous_uri | hierarchy |  
--- | --- | --- | --- | --- |
org/jLAFKha2nDM6kXMKnjv7s3 | 61035901-9824-4186-8ed7-15122daf3760 | org| "" | "{""61035901-9824-4186-8ed7-15122daf3760"": {""type"": ""org"", ""order"": 0}}",
project/qALzZQ6X6nb5A7BAGuPDgW | 0112412c-0747-463a-aa69-573f52d43037 | project | org/mBrhw3qwivJPVuyZUdw6fd | "{""0112412c-0747-463a-aa69-573f52d43037"": {""type"": ""project"", ""order"": 1}, ""e140b44d-81f9-493b-a44e-8d9ada42f842"": {""type"": ""org"", ""order"": 0}}",
folder/KhkJyL8AnCsDuvJpdPDhhH | be07b199-6517-4e89-9c12-026f2a0f2523 | folder | org/mBrhw3qwivJPVuyZUdw6fd | "{""be07b199-6517-4e89-9c12-026f2a0f2523"": {""type"": ""folder"", ""order"": 1}, ""e140b44d-81f9-493b-a44e-8d9ada42f842"": {""type"": ""org"", ""order"": 0}}",

To generate the short uuid, we used the full path to the resource translated as uuids. Eg for the project's uri which
is `qALzZQ6X6nb5A7BAGuPDgW`, we used the
namespace : `61035901-9824-4186-8ed7-15122daf3760/0112412c-0747-463a-aa69-573f52d43037`, I don't know if it makes any
difference from just making a short uuid from `New()`.

#### Hierarchy

The Hierarchy field has the following format

```
{'org/{alias ID}', '{resource_uri}/{resource alias ID}'}
```

For example

```
{'org/mBrhw3qwivJPVuyZUdw6fd', 'repo/UPDpVhNhnWrNLNujxmsz8S'}
```

## API new contract style

Another change of the prototype is the API's json exchange style. As a suggestion is to stop returning UUID and return
complete URL (scheme).

We noticed so far that, we are storing the URI to the resource table and the onlt thing we have to do is to prepend the
server address using the config file.

For example, let us assume the following HTTP request
(_for simplicity we will not require a JWT token and of course skip permissions_) :

`GET |  http://localhost:8080/org/mBrhw3qwivJPVuyZUdw6fd`

```json
{
  "label": "Playground",
  "history": {
    "createdAt": 1610310398125,
    "updatedAt": 1610310398125,
    "createdBy": "Pavlos",
    "updatedBy": "Pavlos"
  },
  "path": {
    "selfUrl": "http://localhost:8080/org/mBrhw3qwivJPVuyZUdw6fd",
    "nextURLs": [
      "http://localhost:8080/repo/UPDpVhNhnWrNLNujxmsz8S",
      "http://localhost:8080/repo/GfKqE2BSX7P6jxwGojfuVW",
      "http://localhost:8080/repo/bW6Krst4zp2dTQi5DKiyoH",
      "http://localhost:8080/project/qALzZQ6X6nb5A7BAGuPDgW",
      "http://localhost:8080/folder/KhkJyL8AnCsDuvJpdPDhhH"
    ]
  }
}
```

What we have achieved so far, with a simple click we can navigate precisely within the platform space to the resources that are attached above and below. 

_Notice_ : of course we can improve the JSON field `nextURLs` to look like this: 
```json
{
  "nextURLs" : {
    "http://localhost:8080/repo/UPDpVhNhnWrNLNujxmsz8S":  {
      "label": "Falkirk Assets",
      "createdBy": "Pavlos"
    },
    "http://localhost:8080/repo/GfKqE2BSX7P6jxwGojfuVW":  {
      "label": "Edinburgh Assets",
      "createdBy": "Pavlos"
    },
    "http://localhost:8080/repo/bW6Krst4zp2dTQi5DKiyoH":  {
      "label": "Glasgow Assets",
      "createdBy": "Pavlos"
    },
    "http://localhost:8080/project/qALzZQ6X6nb5A7BAGuPDgW":  {
      "label": "Tesco IT Department",
      "createdBy": "Az"
    },
    "http://localhost:8080/folder/KhkJyL8AnCsDuvJpdPDhhH":  {
      "label": "Quarter Dataflow Cost",
      "createdBy": "Mark"
    }
  }
}
```

## Future Work

Apply Security on the Resource Path using either the ACL or the RBAC policy framework. 

__Important__ 

This time we should apply policies at the `path URI` instead of the `resource_id`.  
I believe that this small change provides more context to the policy and also we can interrogate if the User has access to the full hierarchy. 