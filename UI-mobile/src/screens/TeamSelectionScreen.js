import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, Alert, Dimensions } from 'react-native';

const { width } = Dimensions.get('window');

export default function TeamSelectionScreen({ navigation }) {
  const [teams, setTeams] = React.useState([]);
  const [selectedTeams, setSelectedTeams] = React.useState([]);

  React.useEffect(() => {
    fetchTeams();
  }, []);

  const fetchTeams = async () => {
    try {
      const API_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080';
      const response = await fetch(`${API_URL}/api/admin/teams`);
      const data = await response.json();
      setTeams(data);
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

  const submitSelection = () => {
    if (selectedTeams.length === 2) {
      navigation.navigate('Results');
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

      <TouchableOpacity style={styles.submitButton} onPress={submitSelection}>
        <Text style={styles.submitButtonText}>Submit Prediction</Text>
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
});
