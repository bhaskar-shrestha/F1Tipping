import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, Alert, Dimensions } from 'react-native';

const { width } = Dimensions.get('window');

export default function TeamSelectionScreen({ navigation, route }) {
  const [teams, setTeams] = React.useState([]);
  const [selectedTeams, setSelectedTeams] = React.useState([]);
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  React.useEffect(() => {
    fetchTeams();
  }, []);

  const fetchTeams = async () => {
    try {
      const API_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_URL}/api/admin/teams`);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      const data = await response.json();
      setTeams(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Error loading teams:', error);
      // Fallback mock data
      setTeams([
        { id: 't1', constructor_name: 'Red Bull' },
        { id: 't2', constructor_name: 'Ferrari' },
        { id: 't3', constructor_name: 'Mercedes' },
        { id: 't4', constructor_name: 'McLaren' },
        { id: 't5', constructor_name: 'Aston Martin' },
      ]);
    }
  };

  const toggleTeam = (teamId) => {
    if (selectedTeams.length < 2) {
      if (selectedTeams.includes(teamId)) {
        setSelectedTeams(selectedTeams.filter(id => id !== teamId));
      } else {
        setSelectedTeams([...selectedTeams, teamId]);
      }
    } else {
      Alert.alert('Limit Reached', 'You can only select 2 teams');
    }
  };

  const submitPrediction = async () => {
    if (selectedTeams.length !== 2) {
      Alert.alert('Error', 'Please select exactly 2 teams');
      return;
    }

    // Get driver IDs from route params (passed from DriverSelectionScreen)
    const selectedDriverIds = route.params?.selectedDrivers || [];
    if (selectedDriverIds.length !== 5) {
      Alert.alert('Error', 'Invalid driver selection');
      return;
    }

    setIsSubmitting(true);
    try {
      const API_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_URL}/api/predictions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 'my-user',
          driver_ids: selectedDriverIds,
          team_ids: selectedTeams,
        }),
      });

      if (!response.ok) {
        // Handle different error statuses
        if (response.status === 400) {
          const errorText = await response.text();
          Alert.alert('Validation Error', errorText || 'Invalid input. Please check your selections.');
          setIsSubmitting(false);
          return;
        } else if (response.status === 404) {
          Alert.alert('Error', 'Server endpoint not found. Please try again later.');
          setIsSubmitting(false);
          return;
        } else {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
      }

      // Success - navigate to Results
      Alert.alert('Success', 'Prediction submitted successfully!');
      navigation.navigate('Results');
    } catch (error) {
      console.error('Error submitting prediction:', error);
      Alert.alert('Error', `Failed to submit prediction: ${error.message}`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Select 2 Teams</Text>

      <ScrollView style={styles.teamList}>
        {teams.map(team => (
          <TouchableOpacity
            key={team.id}
            style={[styles.teamItem, selectedTeams.includes(team.id) && styles.selectedItem]}
            onPress={() => toggleTeam(team.id)}
          >
            <Text style={styles.teamName}>{team.constructor_name}</Text>
            {selectedTeams.includes(team.id) && (
              <Text style={styles.checkmark}>✓</Text>
            )}
          </TouchableOpacity>
        ))}
      </ScrollView>

      <TouchableOpacity 
        style={[styles.submitButton, isSubmitting && styles.submitButtonDisabled]} 
        onPress={submitPrediction}
        disabled={isSubmitting}
      >
        <Text style={styles.submitButtonText}>
          {isSubmitting ? 'Submitting...' : 'Submit Prediction'}
        </Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
    backgroundColor: '#fff',
  },
  header: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  teamList: {
    flex: 1,
    gap: 8,
  },
  teamItem: {
    padding: 12,
    borderRadius: 8,
    backgroundColor: '#f5f5f5',
    flexDirection: 'row',
    alignItems: 'center',
  },
  selectedItem: {
    backgroundColor: '#1976d2',
  },
  teamName: {
    fontSize: 16,
    fontWeight: '600',
  },
  constructorId: {
    fontSize: 14,
    color: '#666',
  },
  checkmark: {
    color: '#fff',
    fontSize: 24,
    marginLeft: 'auto',
  },
  submitButton: {
    padding: 16,
    backgroundColor: '#1976d2',
    borderRadius: 8,
    alignItems: 'center',
  },
  submitButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: '600',
  },
  submitButtonDisabled: {
    opacity: 0.6,
  },
});
