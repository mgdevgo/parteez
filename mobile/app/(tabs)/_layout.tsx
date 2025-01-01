import React from 'react';

import { Tabs } from 'expo-router';

import { BlurView } from 'expo-blur';
import { View, Text, StyleSheet } from 'react-native';
import { SafeAreaView, useSafeAreaInsets } from 'react-native-safe-area-context';

import Ionicicons from '@expo/vector-icons/Ionicons';
import { MapIcon, StarIcon, MagnifyingGlassIcon, RectangleStackIcon } from 'react-native-heroicons/outline';


function TabBarIcon(props: {
  name: React.ComponentProps<typeof Ionicicons>['name'];
  color: string;
}) {
  return <Ionicicons size={28} {...props} />;
}

export default function TabLayout() {
  // const { top } = useSafeAreaInsets();
  // console.log(top)
  return (
    <Tabs screenOptions={{
      headerTransparent: true,
      headerBackground: function() {
        return (
          <BlurView
            intensity={90}
            tint="dark"
            style={{
              height: 96,
              top: 0,
              left: 0,
            }}
          />);
      },
      headerTitleStyle: {
        color: '#ffffff',
      },
      tabBarStyle: {
        position: 'absolute',
        height: 83,
        borderTopColor: 'rgba(234, 234, 234, 0.3)',
      },
      tabBarLabelStyle: {
        fontSize: 13,
        fontWeight: '600',
      },
      tabBarActiveTintColor: '#EA162F',
      tabBarBackground: () => <BlurView intensity={90} tint="dark" style={{ height: 96, bottom: 0, left: 0 }} />,
    }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Тусы',
          tabBarIcon: ({ color }) => <TabBarIcon name="albums" color={color} />,
        }}
      />
      <Tabs.Screen
        name="events"
        options={{
          title: 'Карта',
          tabBarIcon: ({ color }) => <MapIcon size={27} color={color} strokeWidth={1.5} />,
        }}
      />
      <Tabs.Screen
        name="favorite"
        options={{
          title: 'Избранное',
          tabBarIcon: ({ color }) => <StarIcon size={27} color={color} strokeWidth={1.5} />,
        }}
      />
      {/*<Tabs.Screen name="search"*/}
      {/*    options={{*/}
      {/*        title: "Поиск",*/}
      {/*        tabBarIcon: function (tabInfo) {*/}
      {/*            return (*/}
      {/*                <MagnifyingGlassIcon*/}
      {/*                    size={27}*/}
      {/*                    color={tabInfo.color}*/}
      {/*                    strokeWidth={1.5}*/}
      {/*                />*/}
      {/*            );*/}
      {/*        },*/}
      {/*    }}*/}
      {/*/>*/}


    </Tabs>
  );
}