---
title: Promtp
tags:
    - prompt
    - ai
    - golang
---

# Prompt

## Fråga 1

Bygg en golang app som slår samman delar av json-filer som specificeras i ./layer_styles/def/stylefiles.csv

stylefiles.csv:

```csv
aktiv,src_file,target_file,order
1,geonote_basvag_style,springfield_geonote_style,2
1,geonote_hansyn_style,springfield_geonote_style,1
1,geonote_avdelning_style,springfield_geonote_style,0
1,geonote_ovrigt_style,springfield_geonote_style,3
1,arter_style,springfield_geodatapackage_style,50
1,bakgrundskarta_style,springfield_geodatapackage_style,0
1,skifte_style,springfield_geodatapackage_style,1
1,avdelning_style,springfield_geodatapackage_style,2
1,atgard_style,springfield_geodatapackage_style,3
1,anmalan_style,springfield_geodatapackage_style,3.5
1,restriktion_style,springfield_geodatapackage_style,4
1,fornlamning_style,springfield_geodatapackage_style,4.2
1,kulturminne_style,springfield_geodatapackage_style,4.5
1,hansyn_style,springfield_geodatapackage_style,5
0,vag_style,springfield_geodatapackage_style,6
1,drivning_style,springfield_geodatapackage_style,7
1,text_style,springfield_geodatapackage_style,99
```

Exempel json geonote_basvag_style.json:

```json
{
  "$schema": "../spring_style_schema.json",
  "id": "0b11a362d0e44047a2ef65a3850fc867",
  "name": "Mall",
  "version": 8,
  "metadata": {
    "maputnik:renderer": "mbgljs",
    "springGroups": [
      {
        "id": "geonote_basvag",
        "name": "Geonote basväg",
        "collapsed": true,
        "visible": true,
        "springLayers": [
          {
            "id": "geonote_basvag",
            "name": "Basväg",
            "visible": true
          },
          {
            "id": "geonote_avlagg_line",
            "name": "Avlägg",
            "visible": true
          }
        ]
      }
    ],
    "springStyleOrder": 1,
    "springStyleVersion": "0.1"
  },
  "zoom": 0.861983335785597,
  "pitch": 0,
  "center": [
    17.3145660472,
    62.91542224
  ],
  "sprite": "http://localhost:4200/assets/springfield/geonote/sprite/spring_sprites",
  "glyphs": "http://localhost:4200/assets/springfield/geonote/font-glyphs/glyphs/{fontstack}/{range}.pbf",
  "bearing": 0,
  "sources": {
    "geonotes": {
      "type": "geojson",
      "data": "http://localhost:4200/assets/springfield/geonote/geonotes.geojson"
    }
  },
  "layers": [
    {
      "id": "geonote_basvag_background_l",
      "type": "line",
      "metadata": {
        "maputnik:comment": "geonote_basvag urn:sveaskog:atgplan:drivningsinfotyp:basvag",
        "springLayer": "geonote_basvag"
      },
      "source": "geonotes",
      "minzoom": 13,
      "layout": {
        "visibility": "visible"
      },
      "paint": {
        "line-color": "rgba(250, 250, 0, 1)",
        "line-opacity": 1,
        "line-width": 2,
        "line-offset": 4
      },
      "filter": [
        "all",
        [
          "match",
          [
            "geometry-type"
          ],
          [
            "LineString",
            "MultiLineString"
          ],
          true,
          false
        ],
        [
          "==",
          [
            "get",
            "group"
          ],
          "urn:sveaskog:atgplan:drivningsinfotyp:basvag"
        ]
      ]
    },
    {
        "id": "geonote_basvag_background_r",
        "type": "line",
        "metadata": {
          "maputnik:comment": "geonote_basvag urn:sveaskog:atgplan:drivningsinfotyp:basvag",
          "springLayer": "geonote_basvag"
        },
        "source": "geonotes",
        "minzoom": 13,
        "layout": {
          "visibility": "visible"
        },
        "paint": {
          "line-color": "rgba(250, 250, 0, 1)",
          "line-opacity": 1,
          "line-width": 2,
          "line-offset": -4
        },
        "filter": [
          "all",
          [
            "match",
            [
              "geometry-type"
            ],
            [
              "LineString",
              "MultiLineString"
            ],
            true,
            false
          ],
          [
            "==",
            [
              "get",
              "group"
            ],
            "urn:sveaskog:atgplan:drivningsinfotyp:basvag"
          ]
        ]
      }
  ]
}
```

De nycklar i json som skall sammanfogas är "layers", "sources" och "springGroups".
vis sammanslagningen skall objekten innom "layers", "sources" sorteras stigande enligt order i stylefiles.csv medan "springGroups" skall sorteras fallande.