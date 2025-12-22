"""
Training script for ML models
This generates synthetic training data and trains the models
"""
import numpy as np
import sys
import os

# Add parent directory to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from models.risk_forecaster import risk_forecaster, apy_predictor


def generate_synthetic_risk_data(n_samples=1000):
    """Generate synthetic training data for risk forecasting"""
    np.random.seed(42)
    
    # Generate features
    health_factors = np.random.uniform(0.8, 3.0, n_samples)
    collateral_ratios = np.random.uniform(0.3, 0.95, n_samples)
    debt_ratios = 1 - collateral_ratios
    apys = np.random.uniform(2.0, 15.0, n_samples)
    
    X = np.column_stack([health_factors, collateral_ratios, debt_ratios, apys])
    
    # Generate target (risk) based on heuristics
    # Lower health factor = higher risk
    # Higher debt ratio = higher risk
    risk = np.clip(
        1.0 - (health_factors - 0.8) / 2.2 + debt_ratios * 0.3 + np.random.normal(0, 0.1, n_samples),
        0.0, 1.0
    )
    
    return X, risk


def generate_synthetic_apy_data(n_sequences=500, sequence_length=30):
    """Generate synthetic APY time series data"""
    np.random.seed(42)
    
    X = []
    y = []
    
    for _ in range(n_sequences):
        # Generate a time series with trend
        base_apy = np.random.uniform(3.0, 12.0)
        trend = np.random.choice([-0.1, 0, 0.1])
        noise = np.random.normal(0, 0.5, sequence_length)
        
        sequence = base_apy + np.arange(sequence_length) * trend + noise
        sequence = np.clip(sequence, 0.5, 20.0)  # Clamp to reasonable range
        
        # Use last 7 days as features, predict next day
        X.append(sequence[-7:])
        y.append(sequence[-1] + trend + np.random.normal(0, 0.3))
    
    return np.array(X), np.array(y)


def main():
    print("Generating synthetic training data...")
    
    # Train risk forecaster
    print("Training risk forecaster...")
    X_risk, y_risk = generate_synthetic_risk_data(1000)
    risk_forecaster.train(X_risk, y_risk)
    
    # Save risk forecaster
    model_dir = os.path.join(os.path.dirname(__file__), '..', 'models')
    os.makedirs(model_dir, exist_ok=True)
    risk_forecaster.save(os.path.join(model_dir, 'risk_forecaster.pkl'))
    print("Risk forecaster trained and saved!")
    
    # Train APY predictor
    print("Training APY trend predictor...")
    X_apy, y_apy = generate_synthetic_apy_data(500, 30)
    apy_predictor.train(X_apy, y_apy)
    
    # Save APY predictor
    apy_predictor.save(os.path.join(model_dir, 'apy_predictor.pkl'))
    print("APY predictor trained and saved!")
    
    print("\nModel training completed successfully!")
    print(f"Models saved to: {model_dir}")


if __name__ == "__main__":
    main()

