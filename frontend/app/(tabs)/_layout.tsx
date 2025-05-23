import { Tabs } from 'expo-router';
import { MaterialIcons, FontAwesome5, MaterialCommunityIcons } from '@expo/vector-icons';
import { useColorScheme, StyleSheet } from 'react-native';
import { BlurView } from 'expo-blur';

export default function TabLayout() {
  const colorScheme = useColorScheme();
  const isDark = colorScheme === 'dark';

  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarStyle: {
          position: 'absolute',
          borderTopWidth: 0,
          backgroundColor: 'transparent',
          elevation: 0,
        },
        tabBarBackground: () => (
          <BlurView 
            tint={isDark ? 'dark' : 'light'} 
            intensity={80}
            style={StyleSheet.absoluteFill}
          />
        ),
        tabBarInactiveTintColor: isDark ? '#9CA3AF' : '#6B7280',
        tabBarActiveTintColor: isDark ? '#818CF8' : '#6366F1',
        tabBarLabelStyle: {
          fontFamily: 'Inter-Medium',
          fontSize: 12,
          marginBottom: 4,
        },
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
          tabBarIcon: ({ focused, color, size }) => (
            <MaterialIcons 
              name="home" 
              size={size} 
              color={focused ? (isDark ? '#818CF8' : '#6366F1') : (isDark ? '#9CA3AF' : '#6B7280')} 
            />
          ),
        }}
      />
      <Tabs.Screen
        name="games"
        options={{
          title: 'Games',
          tabBarIcon: ({ focused, color, size }) => (
            <MaterialCommunityIcons 
              name="gamepad-variant" 
              size={size} 
              color={focused ? (isDark ? '#818CF8' : '#6366F1') : (isDark ? '#9CA3AF' : '#6B7280')} 
            />
          ),
        }}
      />
      <Tabs.Screen
        name="leaderboard"
        options={{
          title: 'Leaderboard',
          tabBarIcon: ({ focused, color, size }) => (
            <MaterialIcons 
              name="leaderboard" 
              size={size} 
              color={focused ? (isDark ? '#818CF8' : '#6366F1') : (isDark ? '#9CA3AF' : '#6B7280')} 
            />
          ),
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Profile',
          tabBarIcon: ({ focused, color, size }) => (
            <FontAwesome5 
              name="user" 
              size={size - 2} 
              color={focused ? (isDark ? '#818CF8' : '#6366F1') : (isDark ? '#9CA3AF' : '#6B7280')} 
            />
          ),
        }}
      />
      <Tabs.Screen
        name="explore"
        options={{
          title: 'Explore',
          tabBarIcon: ({ focused, color, size }) => (
            <MaterialIcons 
              name="explore" 
              size={size} 
              color={focused ? (isDark ? '#818CF8' : '#6366F1') : (isDark ? '#9CA3AF' : '#6B7280')} 
            />
          ),
        }}
      />
    </Tabs>
  );
}
