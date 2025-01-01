import { useLocalSearchParams } from "expo-router";
import { Text, View } from "react-native";

export default function LoginModalScreen() {
    const { slug } = useLocalSearchParams();
    return (
        <View >
            <Text>Login</Text>
        </View>
    );
}