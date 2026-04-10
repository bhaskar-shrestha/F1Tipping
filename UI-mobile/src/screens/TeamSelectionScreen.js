import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, Alert, Dimensions } from 'react-native';

const { width } = Dimensions.get('window');

export default function TeamSelectionScreen({ navigation }) {
  const [teams, setTeams] = React.useState([]);
  const [selectedTeams, setSelectedTeams] = React.useState([]);

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

  const teamsList = [
    { id: 't1', name: 'Red Bull', constructorId: 'c1' },
    { id: 't2', name: 'Ferrari', constructorId: 'c2' },
    { id: 't3', name: 'Mercedes', constructorId: 'c3' },
    { id: 't4', name: 'McLaren', constructorId: 'c4' },
    { id: 't5', name: 'Aston Martin', constructorId: 'c5' },
  ];

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Select 2 Teams</Text>

      <ScrollView style={styles.teamList}>
        {teamsList.map(team => (
          <TouchableOpacity
            key={team.id}
            style={[styles.teamItem, selectedTeams.includes(team.id) && styles.selectedItem]}
            onPress={() => toggleTeam(team.id)}
          >
            <Text style={styles.teamName}>{team.name}</Text>
            <Text style={styles.constructorId}>{team.constructorId}</Text>
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
