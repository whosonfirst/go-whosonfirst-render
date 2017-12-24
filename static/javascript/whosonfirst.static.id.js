var whosonfirst = whosonfirst || {};
whosonfirst.static = whosonfirst.static || {};

whosonfirst.static.id = (function(){

	var map;
	
	var self = {

		'init': function() {

			self.init_endpoints();
			self.init_map();
			self.init_properties();
			self.init_names();
		},

		'init_endpoints': function() {

			var body = document.body;
			var root = body.getAttribute("data-whosonfirst-uri-endpoint");
			
			if (root){
				whosonfirst.uri.endpoint(root);
			}			
		},
		
		'init_names': function() {
			
			if (typeof(whosonfirst.namify) == 'object'){
				whosonfirst.namify.namify_whosonfirst();
			}
		},

		'init_map': function() {
			
			var geom = document.getElementById("whosonfirst-place");

			var lat = geom.getAttribute("data-latitude");
			var lon = geom.getAttribute("data-longitude");
			
			var centroid_geom = { "type": "Point", "coordinates": [ lon, lat ] };
			var centroid_props = {};
			
			var centroid = {
				"type": "Feature",
				"geometry": centroid_geom,
				"properties": centroid_props
			};
			
			var minlat = geom.getAttribute("data-min-latitude");
			var minlon = geom.getAttribute("data-min-longitude");
			
			var maxlat = geom.getAttribute("data-max-latitude");
			var maxlon = geom.getAttribute("data-max-longitude");	
			
			minlat = parseFloat(minlat);
			minlon = parseFloat(minlon);
			maxlat = parseFloat(maxlat);
			maxlon = parseFloat(maxlon);	
			
			var sw = [ minlat, minlon ];
			var ne = [ maxlat, maxlon ];
			
			var bounds = [ sw, ne ];
			
			var body = document.body;
			var api_key = body.getAttribute("data-mapzen-apikey");
			
			L.Mapzen.apiKey = api_key;
			
			var map_opts = { tangramOptions: {
				scene: L.Mapzen.BasemapStyles.Refill
			}};
			
			map = L.Mapzen.map('map', map_opts);
			
			if ((minlat == maxlat) && (minlon == maxlon)){
				map.setView(sw, 15);
			}
			
			else {
				var fit_opts = {padding: [50, 50]};
				map.fitBounds(bounds, fit_opts);
			}
			
		},

		'init_properties': function(){
			
			var place = document.getElementById("whosonfirst-place");
			var id = place.getAttribute("data-whosonfirst-id");
			
			var data_url = whosonfirst.uri.id2abspath(id)
			// console.log("FETCH", data_url);
			
			var on_success = function(feature){

				var geom = document.getElementById("whosonfirst-place");
				
				var lat = geom.getAttribute("data-latitude");
				var lon = geom.getAttribute("data-longitude");
				
				var centroid_geom = { "type": "Point", "coordinates": [ lon, lat ] };
				var centroid_props = {};
				
				var centroid = {
					"type": "Feature",
					"geometry": centroid_geom,
					"properties": centroid_props
				};
				
				var feature_layer = L.geoJSON(feature);
				feature_layer.addTo(map);
				
				var centroid_layer = L.geoJSON(centroid);
				centroid_layer.addTo(map);

				// sudo put all of this in a wapper function somewhere...
				// (20171224/thisisaaronland)
				
				var props = feature["properties"];
				var props_str = JSON.stringify(props, null, "\t");
				
				var props_raw = document.createElement("pre");
				props_raw.appendChild(document.createTextNode(props_str));
				
				var raw_el = document.getElementById("whosonfirst-properties-raw");
				var pretty_el = document.getElementById("whosonfirst-properties-pretty");

				var button_raw = document.createElement("button");
				button_raw.appendChild(document.createTextNode("show pretty"));
				
				button_raw.onclick = function(){
					raw_el.style.display = "none";
					pretty_el.style.display = "block";					
				};

				var button_pretty = document.createElement("button");
				button_pretty.appendChild(document.createTextNode("show raw"));
				
				button_pretty.onclick = function(){
					pretty_el.style.display = "none";
					raw_el.style.display = "block";					
				};
				
				try {
					var props_pretty = whosonfirst.properties.render(props);

					if (props_pretty){
						raw_el.style.display = "none";
						
						raw_el.appendChild(button_raw);
						raw_el.appendChild(props_raw);
						
						pretty_el.appendChild(button_pretty);
						pretty_el.appendChild(props_pretty);
					}

					else {
						throw "Failed to generate pretty properties";
					}
				}

				catch (e) {
					raw_el.appendChild(button_raw);				
					raw_el.appendChild(props_raw);
					console.log("PRETTY", "ERR", e);
				}

				var hier = props["wof:hierarchy"];
				var hier_pretty = whosonfirst.properties.render(hier);

				var hier_el = document.getElementById("whosonfirst-hierarchy");
				hier_el.appendChild(hier_pretty);

				self.init_names();				
			};
			
			var on_error = function(rsp){
				
				console.log("ERROR", rsp);
				
				var centroid_layer = L.geoJSON(centroid);
				centroid_layer.addTo(map);				
			};
			
			whosonfirst.net.fetch(data_url, on_success, on_error);
		}
	}

	return self;
})();	
