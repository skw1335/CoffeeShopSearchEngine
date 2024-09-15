import json
from shapely_geojson import dumps, Feature, FeatureCollection
from shapely.geometry import Point

with open('data_harvest/data.json') as f:
    d = json.load(f)
   
deduped_json = []
for feature in d['features']:
    coords = (feature['geometry']['coordinates'])
    properties = feature['properties']
    feature = Feature(Point(coords[0], coords[1]), 
                    {
                      'Name': properties["key"],
                      'Review': properties["Review"],
                      'Ratings': properties["Ratings"],
                      'Address': properties["Address"],
                    })
    if feature not in deduped_json:
        deduped_json.append(feature)

feature_collection = FeatureCollection(deduped_json)
with open('deduped_data.json', 'w') as f:
    f.write(dumps(feature_collection, indent=2))

