# goxmlfilter
GoLang to call a rest api which returns xml, and then filter by tags and return column csv of some elements

This can be a useful utility for downloading millions of xml data items from a rest endpoint that serves xml,
filtering them by simple match or discard rules, extracting only a small subset of elements, and dropping into csv files.

Currently it doesn't specify elements using the full path, just the final element name.
I may add the full path match.

## Licence

The [MIT licence](https://github.com/PendaRed/goxmlfilter/blob/master/LICENSE) is used.

## Versions

0.0.1 First version tested January 2020.
0.0.2 January 2020.
    - Optionally turn off cert authority check.
    - Change tage to be an endswith match of the form element/element

## Basic Auth

If you set both environment valiables:
```
xmlfilt_username
xmlfilt_password
```

Then they will be used to set the basicauth header for the query.

## Example restserver
You should run the EgRestServer.go which will serve up 5 people records in the format

```xml
<EgRecords>
    <people>
        <forename>Jonathan</forename>
        <surname>Gibbons</surname>
        <gender>Male</gender>
        <house>12 A Street</house>
    </people>
    <people>
        <forename>Caladan</forename>
        <surname>Brood</surname>
        <gender>Alien</gender>
        <house>20 A Street</house>
    </people>
</EgRecords>
```

## What the code does
It will invoke your URL from the config file and then look for an element which denotes the start of a new record.

For each element until the next delimiter it loads the element and its cdata, and compares the cdata for equality or inequality 
with the configured matchers.  Nothing like xpath, just simple text match.

Records are either include or discarded in this way.

The elements in the record are then compared to the filter extract in the config, and the element can be renamed as well.
These are the only elements which will be output to csv files, with column headers.

## Example config

Look at the [example.json](https://github.com/PendaRed/goxmlfilter/blob/master/src/jgibbons.com/goxmlfilter/example.json) included in the project.

## Config explained

```
{
  // The rest end point you want to invoke
  "query" : "http://localhost:8080/getdata",
  // Some https certs are self certified and not registered with a certificate authority.
  // If not registered set to true, but make sure you know what this means.
  "ignore_cert_authority" : false,
  // How many output rows per file
  "rows_in_file" : 2,
  // How many output files before this will terminate - note it will stop processing
  "num_files" : 2,
  // If you want to see all the elements and values it reads then set true, else false
  "debug_output" : true,
  // Each xml element will have a first element - this cannot be a path, just the element name.
  // If you have arrays of output data this will be how the code knows the row is finished.
  "first_tag" : "forename",
  // An array of equality operations, ie the element must have the value in order to be included
  "filter_equals" : [
    {
      // The element name as suffix path, ie this will match anthing/anything/people/gender 
      // No namespace will be matched
      "element": "people/gender",
      // The value
      "value": "Male"
    }
  ],
  // Similar to equals above, but checks the value does not equal.
  // Do not have the same elements in both the equals and not equals sections
  "filter_not_equals" : [
    {
      "element": "house",
      "value": "12"
    },
    {
      "element": "house",
      "value": "13"
    }
  ],
  // Only the elements listed below will be included in the output csv files. The as value is the column name in the csv header line.
  "filter_extract" : [
    {"element": "people/forename",          "as": "for"},
    {"element": "people/bar",               "as": "bar"},
    {"element": "people/surname",           "as": "sur"}
  ]
}
```

## Generated output

You can say how many files you want output, containing how many rows each.  
It will generate the files using date, time and filenumber.csv for the file name.

### Example data:
```xml
<EgRecords>
  <people>
    <forename>Jonathan</forename>
    <surname>Gibbons</surname>
    <gender>Male</gender>
    <house>12 A Street</house>
  </people>
  <people>
    <forename>Caladan</forename>
    <surname>Brood</surname>
    <gender>Alien</gender>
    <house>20 A Street</house>
    </people>
  <people>
    <forename>Bob</forename>
    <surname>Builder</surname>
    <gender>Male</gender>
    <house>21 A Street</house>
  </people>
  <people>
    <forename>Freddy</forename>
    <surname>Mouse</surname>
    <gender>Male</gender>
    <house>22 A Street</house>
  </people>
  <people>
    <forename>Father</forename>
    <surname>Christmas</surname>
    <gender>Male</gender>
    <house>North Pole</house>
  </people>
</EgRecords>
```

### Example output files:

20200112_181957_1.csv
```
for,bar,sur
Jonathan,,Gibbons
Bob,,Builder
```

20200112_181957_2.csv
```
for,bar,sur
Freddy,,Mouse
Father,,Christmas
```

### Example stdout

This is with debug turned on.
```
URL [http://localhost:8080/getdata]
IgnoreCertAuthority [false]
RowsInFile [2]
NumFiles [2]
DelimTag [forename]
DebugOutput [true]
Filters:
  people/gender = Male
  house != 12
  house != 13
Extract:
  people/forename AS for
  people/bar AS bar
  people/surname AS sur
Appliction: xmlfilter, Version 0.0.2, by Jonathan Gibbons
[18:19:57.944] Calling [http://localhost:8080/getdata]
Setting basic auth for user [jonathan]
[18:19:57.957] Processing Response of size [606] bytes
[18:19:57.957] Writing to file [20200113_181957_1.csv]
   /EgRecords/people/forename: Jonathan
   /EgRecords/people/surname: Gibbons
   /EgRecords/people/gender: Male
   /EgRecords/people/house: 12 A Street
   /EgRecords/people/forename: Caladan
   /EgRecords/people/surname: Brood
   /EgRecords/people/gender: Alien
   /EgRecords/people/house: 20 A Street
   /EgRecords/people/forename: Bob
   /EgRecords/people/surname: Builder
   /EgRecords/people/gender: Male
   /EgRecords/people/house: 21 A Street
[18:19:57.967] Writing to file [20200113_181957_2.csv]
   /EgRecords/people/forename: Freddy
   /EgRecords/people/surname: Mouse
   /EgRecords/people/gender: Male
   /EgRecords/people/house: 22 A Street
   /EgRecords/people/forename: Father
   /EgRecords/people/surname: Christmas
   /EgRecords/people/gender: Male
   /EgRecords/people/house: North Pole
Processing Complete
```
