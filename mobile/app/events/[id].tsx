import { View, Text } from 'react-native'
import React from 'react'
import { Stack, useLocalSearchParams } from 'expo-router'

export default function EventDetails() {
    const { id } = useLocalSearchParams()
  console.log('Visit event details page: id=', id)
    return (
        <View>
            <Stack.Screen options={{ headerTitle: `Event №${id}` }} />
            <Text style={{color:'white'}}>{`Details for event ${id}`}</Text>
        </View>
    )
}