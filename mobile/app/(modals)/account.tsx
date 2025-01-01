import { useLocalSearchParams } from "expo-router";
import { Text, View } from "react-native";

export default function AccountModalScreen() {
  const { slug } = useLocalSearchParams();
  return (
    <View >
      <Text>Account</Text>
    </View>
  );
}