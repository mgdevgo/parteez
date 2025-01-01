import { View, Text, FlatList, Pressable, StyleSheet, ListRenderItem, TouchableOpacity } from 'react-native';
import React, { useState } from 'react';
import { Link, Stack, Tabs } from 'expo-router';
import moment from 'moment';
import { SafeAreaView } from 'react-native-safe-area-context';
import { BlurView } from 'expo-blur';
import { Image } from 'expo-image';

import { useAssets } from 'expo-asset';

const blurhash = '|rF?hV%2WCj[ayj[a|j[az_NaeWBj@ayfRayfQfQM{M|azj[azf6fQfQfQIpWXofj[ayj[j[fQayWCoeoeaya}j[ayfQa{oLj?j[WVj[ayayj[fQoff7azayj[ayj[j[ayofayayayj[fQj[ayayj[ayfjj[j[ayjuayj[';


export default function TodayScreen() {
  // const x = useBottomTabBarHeight()
  const date = moment().format('dddd, MMMM D');
  const [assets, error] = useAssets([require('../../assets/images/blank-big-rave.png')]);
  const dummyArray = [
    { id: '1', artworkUrl: assets ? [0] : null },
    { id: '2', artworkUrl: assets ? [0] : null },
    { id: '3', artworkUrl: assets ? [0] : null },
    { id: '4', artworkUrl: assets ? [0] : null },
    { id: '5', artworkUrl: assets ? [0] : null },
  ];

  const [data, setData] = useState([]);

  async function fetchData() {
    try {
      const response = await fetch('http://192.168.0.104:8080/api/v1/events');
      const data = await response.json();
      console.log(data);
      setData(data.data);

    } catch (error) {
      console.error(error);
    }
  }

  React.useEffect(() => {
    fetchData();
  }, []);

  return (
    <View>
      <Tabs.Screen
        options={{
          headerTransparent: true,
          header: () => (
            <BlurView
              tint="dark"
              intensity={100}
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                height: 48,
              }}
            />
          ),

        }}
      />

      <FlatList
        style={{ paddingHorizontal: 24, backgroundColor: 'black' }}
        ListHeaderComponent={function() {
          return (
            <SafeAreaView style={{ paddingHorizontal: 20 }}>
              <Text style={{
                marginTop: 20,
                textTransform: 'uppercase',
                color: 'white',
                opacity: 0.5,
                fontWeight: '600',
              }}>
                {date}
              </Text>
              <Text style={{ marginTop: 8, color: 'white', fontWeight: '700', fontSize: 36 }}>Today</Text>
            </SafeAreaView>
          );
        }}
        ItemSeparatorComponent={() => <View style={{ height: 20 }} />}
        keyExtractor={(element, i) => element.id}
        data={data}
        renderItem={renderItem}
      >
      </FlatList>

    </View>
  );
}

const renderItem: ListRenderItem<any> = ({ item }) => {
  const eventArtworkPlaceholder = require('../../assets/images/blank-big-rave.png');
  const locationArtworkPlaceholder = require('../../assets/images/favicon.png');
  return (
    <Link
      // href='/modal'
      // href={`/events/${item.id}`}
      href={{ pathname: '/events/[id]', params: { id: item.id } }}
      asChild
    >
      <TouchableOpacity>
        <View
          style={{
            position: 'relative',
            width: '100%',
            borderWidth: 1,
            borderColor: 'red',
            borderRadius: 16,
            justifyContent: 'center',
            alignItems: 'center',
            overflow: 'hidden',
          }}
        >
          <Image
            // source={{ uri: 'https://imgproxy.ra.co/_/quality:66/w:1500/rt:fill/aHR0cHM6Ly9pbWFnZXMucmEuY28vYzEwNGZjNWUzOWY5YmIzNjNiOWY4NDZiM2NhMTc4YWYzNDExOWI0My5wbmc=' }}
            // source={{ uri: require('../../assets/images/blank-big-rave.png') }}
            source={item.artworkUrl}
            placeholder={eventArtworkPlaceholder}
            contentFit="contain"
            transition={1000}
            style={{
              width: '100%',
              height: 375,
            }}
          />

          <BlurView
            style={{
              position: 'absolute',
              bottom: 0,
              left: 0,
              height: 64,
              width: '100%',
              padding: 8,
              display: 'flex',
              flexDirection: 'row',
              alignItems: 'center',
            }}
          >
            <Image
              source={locationArtworkPlaceholder}
              style={{ width: 46, height: 46, borderRadius: 8 }}
              contentFit="contain"
            />
            <View style={{ display: 'flex', flexDirection: 'row', justifyContent: 'center', marginLeft: 3 }}>
              <Text style={{ color: 'white', fontWeight: '500', fontSize: 20 }}>{item.name}</Text>
              <Text style={{ opacity: 60 }}>{`Place: ${item.placeId}`}</Text>
            </View>
          </BlurView>
        </View>
      </TouchableOpacity>
    </Link>);
};