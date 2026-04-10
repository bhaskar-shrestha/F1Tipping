import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, Alert, Dimensions } from 'react-native';
import { Link } from 'react-navigation';

const { width } = Dimensions.get('window');

export default function DriverSelectionScreen({ navigation }) {
  const [drivers, setDrivers] = useState([]);
  const [selectedDrivers, setSelectedDrivers] = useState([]);

  useEffect(() => {
    // Load drivers from API
    fetchDrivers();
  }, []);

  const fetchDrivers = async () => {
    try {
      // In real app: API.get('/api/admin/drivers')
      setDrivers([
        { id: 'd1', name: 'Max Verstappen', constructorName: 'Red Bull' },
        { id: 'd2', name: 'Sergio Perez', constructorName: 'Red Bull' },
        { id: 'd3', name: 'Charles Leclerc', constructorName: 'Ferrari' },
        { id: 'd4', name: 'Carlos Sainz', constructorName: 'Ferrari' },
        { id: 'd5', name: 'Lewis Hamilton', constructorName: 'Mercedes' },
        { id: 'd6', name: 'George Russell', constructorName: 'Mercedes' },
        { id: 'd7', name: 'Lando Norris', constructorName: 'McLaren' },
        { id: 'd8', name: 'Oscar Piastri', constructorName: 'McLaren' },
      ]);
    } catch (error) {
      console.error('Error loading drivers:', error);
    }
  };

  const toggleDriver = (driverId) => {
    if (selectedDrivers.length < 5) {
      if (selectedDrivers.includes(driverId)) {
        setSelectedDrivers(selectedDrivers.filter(id => id !== driverId));
      } else {
        setSelectedDrivers([...selectedDrivers, driverId]);
      }
    } else {
      Alert.alert('Limit Reached', 'You can only select 5 drivers');
    }
  };

  const submitSelection = () => {
    if (selectedDrivers.length === 5) {
      navigation.navigate('TeamSelection');
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.header}>Select 5 Drivers</Text>

      <ScrollView style={styles.driverList}>
        {drivers.map(driver => (
          <TouchableOpacity
            key={driver.id}
            style={[styles.driverItem, selectedDrivers.includes(driver.id) && styles.selectedItem]}
            onPress={() => toggleDriver(driver.id)}
          >
            <Text style={styles.driverName}>{driver.name}</Text>
            <Text style={styles.constructorName}>{driver.constructorName}</Text>
            {selectedDrivers.includes(driver.id) && (
              <Text style={styles.checkmark}>✓</Text>
            )}
          </TouchableOpacity>
        ))}
      </ScrollView>

      <TouchableOpacity style={styles.submitButton} onPress={submitSelection}>
        <Text style={styles.submitButtonText}>Next</Text>
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
  driverList: {
    flex: 1,
    gap: 8,
  },
  driverItem: {
    padding: 12,
    borderRadius: 8,
    backgroundColor: '#f5f5f5',
    flexDirection: 'row',
    alignItems: 'center',
  },
  selectedItem: {
    backgroundColor: '#1976d2',
  },
  driverName: {
    fontSize: 16,
    fontWeight: '600',
  },
  constructorName: {
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
