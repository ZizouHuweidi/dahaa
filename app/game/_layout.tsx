import { Stack } from 'expo-router';

export default function GameLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="create"
        options={{
          title: 'Create Game',
        }}
      />
      <Stack.Screen
        name="join"
        options={{
          title: 'Join Game',
        }}
      />
    </Stack>
  );
} 