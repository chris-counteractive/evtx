package evtx

// EVTX is in XML but this is hard for us to query. So we try to
// normalize some common XML patterns into something which is easier
// to work with.

/* Sometimes the XML looks like:

<EventData>
   <Data name="Thing1"> %Subst% </Data>
   <Data name="Thing2"> %Subst% </Data>
   <Data name="Thing3"> %Subst% </Data>
</Eventdata>

We convert it to json like this:

"Eventdata": {
    "Data": [
       {
         "Name": "thing1",
         "": %Subst%
       },
       {
         "Name": "thing2",
         "": %Subst%
       },
       {
         "Name": "thing3",
         "": %Subst%
       },

    ]
}

But this is really hard to work with. To simplify processing we need
to convert it to:

"EventData": {
   "Thing1": %Subst,
   "Thing1": %Subst,
   "Thing1": %Subst,
}
*/

func NormalizeEventData(expanded interface{}) {
	data, ok := expanded.(map[string]interface{})
	if !ok {
		return
	}
	event_data, pres := data["EventData"]
	if !pres {
		return
	}
	event_data_map, ok := event_data.(map[string]interface{})
	if !ok {
		return
	}

	data_tag, pres := event_data_map["Data"]
	if !pres {
		return
	}

	data_array, ok := data_tag.([]interface{})
	if !ok {
		return
	}

	result := make(map[string]interface{})
	for _, item := range data_array {
		item_map, ok := item.(map[string]interface{})
		if !ok {
			return
		}

		// Look for name and "" pairs.
		name_any, pres := item_map["Name"]
		if !pres {
			return
		}

		name, ok := name_any.(string)
		if !ok {
			return
		}

		value, pres := item_map["Value"]
		if !pres {
			return
		}
		result[name] = value
	}

	data["EventData"] = result
}