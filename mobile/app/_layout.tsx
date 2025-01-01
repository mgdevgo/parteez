import * as React from 'react';

import '@/tailwind.css';
import { useFonts } from 'expo-font';
import { Ionicons } from '@expo/vector-icons';
import { TouchableOpacity } from 'react-native';
import { Stack, SplashScreen, useRouter } from 'expo-router';

export default function RootLayout() {
  const [loaded, error] = useFonts({
    // SpaceMono: require("../assets/fonts/SpaceMono-Regular.ttf"),
    ...Ionicons.font,
  });

  // Expo Router uses Error Boundaries to catch errors in the navigation tree.
  React.useEffect(() => {
    if (error) throw error;
  }, [error]);

  React.useEffect(() => {
    if (loaded) {
      SplashScreen.hideAsync();
    }
  }, [loaded]);

  if (!loaded) {
    return null;
  }

  return (
    <RootLayoutNavigation />
  )
}

function RootLayoutNavigation() {
  const router = useRouter();

  return (
    <Stack
      screenOptions={{
        contentStyle: {
          backgroundColor: '#000',
        }
      }}
    >
      <Stack.Screen name="(modals)/login" options={{
        title: 'Login or Create Account',
        presentation: 'modal',
        headerBackTitle: 'Close',
        headerLeft: function() {
          return (
            <TouchableOpacity onPress={() => router.back()}>
              <Ionicons name="close-outline" size={28} />
            </TouchableOpacity>
          )
        },
      }}
      />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="events/[id]" options={{ headerShown: true }} />

    </Stack>
  );
}
