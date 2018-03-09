# pivotal-tracker-resource

An output-only resource for interacting with Pivotal Tracker stories. Heavily inspired by the [jira-resource](https://github.com/danrspencer/jira-resource).

## quick example

See more in `_examples/`

```yaml
jobs:
# creates pivotal tracker stories to upgrade all your environments
#  when new versions of a tile product are available
- name: create-story-new-version-p-mysql
  plan:
  - get: pivnet-product
    resource: p-mysql
    trigger: true
    params:
      globs: []
  - task: generate-file-content
    config:
      platform: linux
      image_resource:
        type: docker-image
        source: {repository: czero/cflinuxfs2}
      inputs:
      - name: pivnet-product
      outputs:
      - name: generate-file-content
      run:
        path: bash
        args:
        - -c
        - |
          desired_version=$(jq --raw-output '.Release.Version' < ./pivnet-product/metadata.json)
          echo "$desired_version" > generate-file-content/output.txt
  - aggregate:
    - put: tracker
      params:
        name: Upgrade to p-mysql v$NAME_FILE in Sandbox
        name_file: generate-file-content/output.txt
        owner_ids:
        - 3016456
        story_type: chore
        labels:
        - name: sandbox
        - name: p-mysql
        - name: tile-upgrade
    - put: tracker
      params:
        name: Upgrade to p-mysql v$NAME_FILE in Dev
        name_file: generate-file-content/output.txt
        owner_ids:
        - 3016456
        story_type: chore
        labels:
        - name: dev
        - name: p-mysql
        - name: tile-upgrade
    - put: tracker
      params:
        name: Upgrade to p-mysql v$NAME_FILE in Prod
        name_file: generate-file-content/output.txt
        owner_ids:
        - 3016456
        story_type: chore
        labels:
        - name: prod
        - name: p-mysql
        - name: tile-upgrade
```

## source configuration

```yaml
resources:
- name: tracker
  type: pivotal-tracker-resource
  source:
    project_id: 12345
    token: abcdefghijklmnop
```

## resource type configuration

```yaml
resource_types:
- name: pivotal-tracker-resource
  type: docker-image
  source:
    repository: aegershman/pivotal-tracker-resource
    tag: latest
```

## behavior

### `out`: creates or updates a tracker story

#### parameters

* `name`: *required*
* `name_file`: *optional* if you reference `$NAME_FILE` in the text string of `name`, you can specify the file whose contents will replace references to `$NAME_FILE`

```yaml
- put: tracker
  params:
    name: This ticket name is static

- put: tracker
  params:
    name: Upgrade to PAS v$NAME_FILE in Sandbox
    name_file: pivnet-product/version
```

WARNING: For some reason, this has a very difficult time matching story names when letters, numbers, or punctuation are merged together. This might be an issue with my code, or maybe I'm misunderstanding something about pivotal tracker's search API. Regardless, until this is resolved, it's recommended you separate your version numbers from any other text. For example:

```yaml
- put: tracker
  params:
    name: This fails to find a match v1.1.10

- put: tracker
  params:
    name: But this works fine 1.1.10
```

* `description`: *optional*

```yaml
- put: tracker
  params:
    name: blah
    description: blah blah blah blah
```

* `owner_ids`: *optional* integer IDs of the story owners. In the future I'd like to have this be a list of emails or userIDs, but for now it has to be the user's database ID

```yaml
- put: tracker
  params:
    name: blah
    owner_ids:
    - 1234
    - 9999
```

* `story_type`: *optional* acceptable values are `feature, bug, chore, release`

```yaml
- put: tracker
  params:
    name: blah
    story_type: chore
```

* `current_state`: *optional* acceptable values are `accepted, delivered, finished, started, rejected, planned, unstarted, unscheduled`

```yaml
- put: tracker
  params:
    name: blah
    current_state: started
```

* `labels`: *optional*

```yaml
- put: tracker
  params:
    name: blah
    labels:
    - name: blah-automation
    - name: blah-dev
```

## contributing

PRs, suggestions, comments, etc., more than welcome. I'm still learning `go`, so if you have any feedback at all, please make them known.