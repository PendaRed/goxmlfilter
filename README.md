# goxmlfilter
GoLang to call a rest api which returns xml, and then filter by tags and return column csv of some elements

This can be a useful utility for downloading millions of xml data items from a rest endpoint that serves xml,
filtering them by simple match or discard rules, extracting only a small subset of elements, and dropping into csv files.

Currently it doesn't specify elements using the full path, just the final element name.
I may add the full path match.

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

## Generated output

You can say how many files you want output, containing how many rows each.  
It will generate the files using date, time and filenumber.csv for the file name.

