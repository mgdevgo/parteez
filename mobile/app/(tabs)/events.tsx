import * as React from 'react';

import { Tabs } from 'expo-router';

import { View, Text, Button, Pressable, TouchableOpacity } from 'react-native';
import MapView from 'react-native-maps';

export default function EventsScreen() {
  return (

    <View>
      <Tabs.Screen
        options={{
          headerRight: function() {
            return (
              <TouchableOpacity>
                <Text style={{ color: '#fff' }}>Список</Text>
              </TouchableOpacity>
            );
          },
        }}
      />
      <Text>Map</Text>
      <MapView style={{ width: '100%', height: '100%' }}
               initialRegion={{
                 latitude: 59.940168,
                 longitude: 30.326154,
                 latitudeDelta: 0.0922,
                 longitudeDelta: 0.0421,
               }}
      />
    </View>
  );
}