import { StyleSheet, View, TouchableOpacity, ScrollView, useColorScheme } from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { Text } from '@/components/ui/Text';

const ActionButton = ({ 
  icon, 
  title, 
  subtitle, 
  onPress,
  color
}: {
  icon: string;
  title: string;
  subtitle: string;
  onPress: () => void;
  color: string;
}) => (
  <TouchableOpacity style={styles.actionButton} onPress={onPress}>
    <View style={[styles.actionIcon, { backgroundColor: `${color}20` }]}>
      <MaterialIcons name={icon as any} size={24} color={color} />
    </View>
    <View style={styles.actionTextContainer}>
      <Text variant="h4" weight="semiBold" style={styles.actionTitle}>
        {title}
      </Text>
      <Text variant="caption" color="secondary" style={styles.actionSubtitle}>
        {subtitle}
      </Text>
    </View>
    <MaterialIcons 
      name="chevron-right" 
      size={24} 
      color={useColorScheme() === 'dark' ? '#9CA3AF' : '#6B7280'} 
    />
  </TouchableOpacity>
);

export default function HomeScreen() {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';

  const handleCreateGame = () => {
    // Navigate to create game screen
    console.log('Create game');
  };

  const handleJoinGame = () => {
    // Navigate to join game screen
    console.log('Join game');
  };

  const handleHowToPlay = () => {
    // Navigate to how to play screen
    console.log('How to play');
  };

  return (
    <ScrollView 
      style={[styles.container, isDark ? styles.darkContainer : styles.lightContainer]}
      contentContainerStyle={styles.contentContainer}
    >
      <View style={styles.header}>
        <Text variant="h1" weight="bold" style={styles.title}>
          Welcome to
        </Text>
        <Text 
          variant="h1" 
          weight="bold" 
          style={[styles.title, { color: isDark ? '#818CF8' : '#6366F1' }]}
        >
          Dahaa
        </Text>
        <Text variant="body" color="secondary" style={styles.subtitle}>
          The ultimate bluffing and deduction game
        </Text>
      </View>

      <View style={styles.actionsContainer}>
        <Text variant="h3" weight="semiBold" style={styles.sectionTitle}>
          Quick Actions
        </Text>
        
        <ActionButton
          icon="add-circle-outline"
          title="Create Game"
          subtitle="Start a new game with friends"
          onPress={handleCreateGame}
          color="#818CF8"
        />
        
        <ActionButton
          icon="group-add"
          title="Join Game"
          subtitle="Join an existing game"
          onPress={handleJoinGame}
          color="#10B981"
        />
        
        <ActionButton
          icon="help-outline"
          title="How to Play"
          subtitle="Learn the rules"
          onPress={handleHowToPlay}
          color="#F59E0B"
        />
      </View>

      <View style={styles.recentGamesContainer}>
        <View style={styles.sectionHeader}>
          <Text variant="h3" weight="semiBold">
            Recent Games
          </Text>
          <TouchableOpacity>
            <Text variant="body" color="primary" weight="medium">
              See All
            </Text>
          </TouchableOpacity>
        </View>
        
        <View style={styles.emptyState}>
          <MaterialIcons 
            name="sports-esports" 
            size={48} 
            color={isDark ? '#4B5563' : '#9CA3AF'} 
          />
          <Text variant="body" color="secondary" style={styles.emptyStateText}>
            No recent games
          </Text>
        </View>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  lightContainer: {
    backgroundColor: '#F9FAFB',
  },
  darkContainer: {
    backgroundColor: '#111827',
  },
  contentContainer: {
    padding: 16,
    paddingTop: 24,
  },
  header: {
    marginBottom: 32,
  },
  title: {
    fontSize: 36,
    lineHeight: 40,
    marginBottom: 4,
  },
  subtitle: {
    fontSize: 16,
    marginTop: 8,
  },
  actionsContainer: {
    marginBottom: 32,
  },
  sectionTitle: {
    marginBottom: 16,
  },
  actionButton: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 16,
    borderRadius: 12,
    backgroundColor: 'transparent',
    marginBottom: 12,
    borderWidth: 1,
    borderColor: 'rgba(0, 0, 0, 0.05)',
  },
  actionIcon: {
    width: 48,
    height: 48,
    borderRadius: 12,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 16,
  },
  actionTextContainer: {
    flex: 1,
  },
  actionTitle: {
    marginBottom: 4,
  },
  actionSubtitle: {
    fontSize: 14,
  },
  recentGamesContainer: {
    marginBottom: 24,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  emptyState: {
    backgroundColor: 'rgba(0, 0, 0, 0.02)',
    borderRadius: 12,
    padding: 32,
    alignItems: 'center',
    justifyContent: 'center',
  },
  emptyStateText: {
    marginTop: 12,
    textAlign: 'center',
  },
});
