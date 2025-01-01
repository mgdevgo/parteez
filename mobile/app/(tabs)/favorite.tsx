import React from 'react';

import { Link } from 'expo-router';

import { View, Text, Pressable } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

export default function Favorite() {
  return (
      <View style={{ paddingTop: 96, backgroundColor: "black", height: '100%'}}>
        <Text style={{color: 'white'}}>Favorite</Text>
        <Link href={{ pathname: '/(modals)/login' }} asChild>
          <Pressable>
            <Text style={{ color: 'black' }}>Login</Text>
          </Pressable>
        </Link>
      </View>
  );
}