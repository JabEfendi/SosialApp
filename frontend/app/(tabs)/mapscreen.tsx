import React from 'react';
import { View, StyleSheet } from 'react-native';
import { WebView } from 'react-native-webview';

const MAPBOX_TOKEN = process.env.EXPO_PUBLIC_MAPBOX_TOKEN;

export default function MapScreen() {
  const html = `
<!DOCTYPE html>
<html>
<head>
  <meta name="viewport" content="initial-scale=1.0, maximum-scale=1.0" />
  <link href="https://api.mapbox.com/mapbox-gl-js/v2.15.0/mapbox-gl.css" rel="stylesheet" />
  <script src="https://api.mapbox.com/mapbox-gl-js/v2.15.0/mapbox-gl.js"></script>
  <style>
    body, html, #map {
      margin: 0;
      padding: 0;
      width: 100%;
      height: 100%;
    }
  </style>
</head>
<body>
  <div id="map"></div>

  <script>
    mapboxgl.accessToken = '${MAPBOX_TOKEN}';

    const map = new mapboxgl.Map({
      container: 'map',
      style: 'mapbox://styles/mapbox/streets-v12',
      center: [106.816666, -6.200000],
      zoom: 14
    });

    let marker = null;

    map.on('click', function(e) {
      const { lng, lat } = e.lngLat;

      if (marker) marker.remove();

      marker = new mapboxgl.Marker()
        .setLngLat([lng, lat])
        .addTo(map);

      window.ReactNativeWebView.postMessage(
        JSON.stringify({ latitude: lat, longitude: lng })
      );
    });
  </script>
</body>
</html>
`;

  return (
    <View style={styles.container}>
      <WebView
        originWhitelist={['*']}
        source={{ html }}
        javaScriptEnabled
        onMessage={(event) => {
          const data = JSON.parse(event.nativeEvent.data);
          console.log('Lokasi dipilih:', data);
        }}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1 }
});