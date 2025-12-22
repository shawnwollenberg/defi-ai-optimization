"""
Risk Forecasting Models for DeFi Positions
"""
import numpy as np
from sklearn.ensemble import RandomForestRegressor
from sklearn.preprocessing import StandardScaler
import pickle
import os


class RiskForecaster:
    """ML model for predicting liquidation risk"""
    
    def __init__(self):
        self.model = None
        self.scaler = StandardScaler()
        self.is_trained = False
        
    def train(self, X, y):
        """
        Train the risk forecasting model
        X: Features (health_factor, collateral_ratio, debt_ratio, etc.)
        y: Target (liquidation_risk: 0-1)
        """
        # Scale features
        X_scaled = self.scaler.fit_transform(X)
        
        # Train Random Forest model
        self.model = RandomForestRegressor(
            n_estimators=100,
            max_depth=10,
            random_state=42,
            n_jobs=-1
        )
        self.model.fit(X_scaled, y)
        self.is_trained = True
        
    def predict(self, features):
        """
        Predict liquidation risk
        features: dict with health_factor, total_collateral, total_debt, etc.
        """
        if not self.is_trained:
            # Return default prediction if model not trained
            return self._default_prediction(features)
        
        # Extract features in order: health_factor, collateral_ratio, debt_ratio, apy
        X = np.array([[
            features.get('health_factor', 1.5),
            features.get('collateral_ratio', 0.8),
            features.get('debt_ratio', 0.5),
            features.get('apy', 5.0),
        ]])
        
        X_scaled = self.scaler.transform(X)
        risk = self.model.predict(X_scaled)[0]
        
        # Clamp to [0, 1]
        risk = max(0.0, min(1.0, risk))
        
        return risk
    
    def _default_prediction(self, features):
        """Default risk prediction based on heuristics"""
        health_factor = features.get('health_factor', 1.5)
        
        # Simple heuristic: lower health factor = higher risk
        if health_factor < 1.1:
            return 0.9  # Very high risk
        elif health_factor < 1.3:
            return 0.6  # High risk
        elif health_factor < 1.5:
            return 0.3  # Medium risk
        else:
            return 0.1  # Low risk
    
    def save(self, filepath):
        """Save the trained model"""
        if self.is_trained:
            with open(filepath, 'wb') as f:
                pickle.dump({
                    'model': self.model,
                    'scaler': self.scaler,
                }, f)
    
    def load(self, filepath):
        """Load a trained model"""
        if os.path.exists(filepath):
            with open(filepath, 'rb') as f:
                data = pickle.load(f)
                self.model = data['model']
                self.scaler = data['scaler']
                self.is_trained = True


class APYTrendPredictor:
    """ML model for predicting APY trends"""
    
    def __init__(self):
        self.model = None
        self.scaler = StandardScaler()
        self.is_trained = False
    
    def train(self, X, y):
        """
        Train the APY trend predictor
        X: Historical APY data (time series features)
        y: Future APY values
        """
        X_scaled = self.scaler.fit_transform(X)
        
        self.model = RandomForestRegressor(
            n_estimators=100,
            max_depth=8,
            random_state=42,
            n_jobs=-1
        )
        self.model.fit(X_scaled, y)
        self.is_trained = True
    
    def predict(self, historical_apy):
        """
        Predict future APY based on historical data
        historical_apy: List of historical APY values
        """
        if not self.is_trained or len(historical_apy) < 7:
            # Return average if model not trained or insufficient data
            return np.mean(historical_apy) if historical_apy else 5.0
        
        # Use last 7 days as features
        features = np.array([historical_apy[-7:]])
        features_scaled = self.scaler.transform(features)
        
        predicted = self.model.predict(features_scaled)[0]
        return max(0.0, predicted)  # APY can't be negative
    
    def predict_trend(self, historical_apy):
        """
        Predict trend direction (increasing, decreasing, stable)
        """
        if len(historical_apy) < 2:
            return "stable"
        
        predicted = self.predict(historical_apy)
        recent_avg = np.mean(historical_apy[-7:]) if len(historical_apy) >= 7 else np.mean(historical_apy)
        
        diff = predicted - recent_avg
        threshold = 0.5  # 0.5% change threshold
        
        if diff > threshold:
            return "increasing"
        elif diff < -threshold:
            return "decreasing"
        else:
            return "stable"
    
    def save(self, filepath):
        """Save the trained model"""
        if self.is_trained:
            with open(filepath, 'wb') as f:
                pickle.dump({
                    'model': self.model,
                    'scaler': self.scaler,
                }, f)
    
    def load(self, filepath):
        """Load a trained model"""
        if os.path.exists(filepath):
            with open(filepath, 'rb') as f:
                data = pickle.load(f)
                self.model = data['model']
                self.scaler = data['scaler']
                self.is_trained = True


# Global model instances
risk_forecaster = RiskForecaster()
apy_predictor = APYTrendPredictor()

# Try to load pre-trained models if they exist
MODEL_DIR = os.path.join(os.path.dirname(__file__), '..', 'models')
os.makedirs(MODEL_DIR, exist_ok=True)

risk_forecaster.load(os.path.join(MODEL_DIR, 'risk_forecaster.pkl'))
apy_predictor.load(os.path.join(MODEL_DIR, 'apy_predictor.pkl'))

