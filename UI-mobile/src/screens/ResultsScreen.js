import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';

export default function ResultsScreen() {
  const [results, setResults] = React.useState([]);

  React.useEffect(() => {
    // In real app: fetch results from API
    setResults([
      {
        id: 'r1',
        drivers: 'Max Verstappen, Sergio Perez, Charles Leclerc, Carlos Sainz, Lewis Hamilton',
        teams: 'Red Bull, Ferrari',
        sprintPoints: 25,
        racePoints: 60,
        totalPoints: 85,
      },
    ]);
  }, []);

  const getPointsColor = (points) => {
    if (points >= 50) return '#4caf50';
    if (points >= 30) return '#ff9800';
    return '#f5f5f5';
  };

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Your Results</Text>

      {results.length === 0 ? (
        <Text style={styles.noResults}>No results yet</Text>
      ) : (
        results.map(result => (
          <View key={result.id} style={styles.card}>
            <Text style={styles.sectionTitle}>Prediction</Text>
            <Text style={styles.driversText}>
              Drivers: {result.drivers}
            </Text>
            <Text style={styles.teamsText}>
              Teams: {result.teams}
            </Text>

            <Text style={styles.sectionTitle}>Points</Text>
            <View style={styles.pointsRow}>
              <View style={styles.pointsBox}>
                <Text style={styles.pointsLabel}>Sprint</Text>
                <Text style={[styles.pointsValue, { color: getPointsColor(result.sprintPoints) }]}>
                  {result.sprintPoints}
                </Text>
              </View>
              <View style={styles.pointsBox}>
                <Text style={styles.pointsLabel}>Race</Text>
                <Text style={[styles.pointsValue, { color: getPointsColor(result.racePoints) }]}>
                  {result.racePoints}
                </Text>
              </View>
            </View>

            <View style={styles.totalPointsBox}>
              <Text style={styles.totalPointsLabel}>Total Points</Text>
              <Text style={[styles.totalPointsValue, { color: getPointsColor(result.totalPoints) }]}>
                {result.totalPoints}
              </Text>
            </View>
          </View>
        ))
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
    backgroundColor: '#f5f5f5',
  },
  header: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 16,
  },
  noResults: {
    fontSize: 16,
    color: '#666',
    textAlign: 'center',
    marginTop: 40,
  },
  card: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 8,
  },
  driversText: {
    fontSize: 14,
    color: '#333',
    marginBottom: 4,
  },
  teamsText: {
    fontSize: 14,
    color: '#666',
    marginBottom: 12,
  },
  sectionTitlePoints: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 8,
  },
  pointsRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginBottom: 12,
  },
  pointsBox: {
    alignItems: 'center',
    padding: 8,
    borderRadius: 8,
    backgroundColor: '#f5f5f5',
  },
  pointsLabel: {
    fontSize: 12,
    color: '#666',
  },
  pointsValue: {
    fontSize: 24,
    fontWeight: 'bold',
  },
  totalPointsBox: {
    backgroundColor: '#1976d2',
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
  },
  totalPointsLabel: {
    color: '#fff',
    fontSize: 14,
  },
  totalPointsValue: {
    color: '#fff',
    fontSize: 32,
    fontWeight: 'bold',
    marginTop: 4,
  },
});
