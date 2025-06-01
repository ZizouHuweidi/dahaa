import * as React from 'react';
import { View } from 'react-native';
import { useRouter } from 'expo-router';
import { Button } from '~/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import { Text } from '~/components/ui/text';
import { useGameStore, Player } from '~/lib/store/game';

export default function GamePlayScreen() {
  const router = useRouter();
  const { currentGame, isLoading, error } = useGameStore();
  const [answer, setAnswer] = React.useState('');
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  // Placeholder for current round/question
  const currentQuestion = currentGame?.rounds?.[0]?.question || 'Waiting for question...';

  const handleSubmit = async () => {
    if (!answer.trim()) return;
    setIsSubmitting(true);
    // TODO: Submit answer to backend
    setTimeout(() => {
      setIsSubmitting(false);
      setAnswer('');
      // Optionally navigate or show feedback
    }, 1000);
  };

  if (!currentGame) {
    return (
      <View className="flex-1 p-6 bg-background justify-center items-center">
        <Text className="text-lg">No active game found.</Text>
        <Button className="mt-4" onPress={() => router.replace('/')}>Return Home</Button>
      </View>
    );
  }

  return (
    <View className="flex-1 p-6 bg-background">
      <Card className="w-full max-w-md mx-auto">
        <CardHeader>
          <CardTitle>Game In Progress</CardTitle>
          <CardDescription>Answer the question below:</CardDescription>
        </CardHeader>
        <CardContent className="gap-4">
          {error && <Text className="text-destructive text-center">{error}</Text>}
          <Text className="text-xl font-bold mb-4">{currentQuestion}</Text>
          <Input
            placeholder="Your answer..."
            value={answer}
            onChangeText={setAnswer}
            editable={!isSubmitting}
          />
        </CardContent>
        <CardFooter className="flex-col gap-2">
          <Button
            className="w-full"
            onPress={handleSubmit}
            disabled={isSubmitting || !answer.trim()}
          >
            <Text className="text-primary-foreground">
              {isSubmitting ? 'Submitting...' : 'Submit Answer'}
            </Text>
          </Button>
          <View className="mt-4">
            <Text className="font-semibold mb-2">Players:</Text>
            {currentGame.players.map((player: Player) => (
              <Text key={player.id} className="text-base">
                {player.name}
                {currentGame.rounds?.[0]?.answer_pool?.fake_answers?.some(
                  (a) => a.player_id === player.id
                ) ? ' (Answered)' : ' (Not Answered)'}
              </Text>
            ))}
          </View>
        </CardFooter>
      </Card>
    </View>
  );
} 