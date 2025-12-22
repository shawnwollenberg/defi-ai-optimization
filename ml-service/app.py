"""
FastAPI ML Service for DeFi Risk Forecasting
"""
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
import os
from dotenv import load_dotenv
from models.risk_forecaster import risk_forecaster, apy_predictor

load_dotenv()

app = FastAPI(title="DeFi ML Service", version="1.0.0")


class RiskForecastRequest(BaseModel):
    """Request model for risk forecasting"""
    user_address: str
    positions: List[dict]
    health_factor: Optional[float] = None
    total_collateral: Optional[float] = None
    total_debt: Optional[float] = None


class APYTrendRequest(BaseModel):
    """Request model for APY trend analysis"""
    protocol: str
    asset: str
    historical_apy: List[float]
    days: int = 30


class RiskForecastResponse(BaseModel):
    """Response model for risk forecasting"""
    liquidation_risk: float
    risk_level: str
    recommendations: List[str]
    confidence: float


class APYTrendResponse(BaseModel):
    """Response model for APY trend analysis"""
    predicted_apy: float
    trend: str
    confidence: float
    recommendation: str


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "ml-service"}


@app.post("/api/v1/risk/forecast", response_model=RiskForecastResponse)
async def forecast_risk(request: RiskForecastRequest):
    """
    Forecast liquidation risk for a user's DeFi positions
    """
    # Calculate risk from positions
    total_collateral = request.total_collateral or 0.0
    total_debt = request.total_debt or 0.0
    health_factor = request.health_factor
    
    # Calculate health factor if not provided
    if health_factor is None and total_debt > 0:
        health_factor = (total_collateral * 0.8) / total_debt  # Simplified calculation
    elif health_factor is None:
        health_factor = 2.0  # Safe default
    
    # Prepare features for ML model
    features = {
        'health_factor': health_factor,
        'collateral_ratio': total_collateral / (total_collateral + total_debt) if (total_collateral + total_debt) > 0 else 1.0,
        'debt_ratio': total_debt / (total_collateral + total_debt) if (total_collateral + total_debt) > 0 else 0.0,
        'apy': 5.0,  # Default APY
    }
    
    # Get average APY from positions if available
    if request.positions:
        avg_apy = sum(p.get('apy', 5.0) for p in request.positions) / len(request.positions)
        features['apy'] = avg_apy
    
    # Predict risk
    risk = risk_forecaster.predict(features)
    
    # Determine risk level
    if risk >= 0.7:
        risk_level = "critical"
        recommendations = [
            "Immediate action required: Reduce debt or add collateral",
            "Risk of liquidation is very high",
            "Consider closing positions"
        ]
    elif risk >= 0.5:
        risk_level = "high"
        recommendations = [
            "Consider reducing leverage",
            "Monitor health factor closely",
            "Add collateral to improve safety margin"
        ]
    elif risk >= 0.3:
        risk_level = "medium"
        recommendations = [
            "Monitor health factor regularly",
            "Consider rebalancing if APY drops"
        ]
    else:
        risk_level = "low"
        recommendations = [
            "Portfolio is in good health",
            "Continue monitoring"
        ]
    
    # Confidence based on model training status
    confidence = 0.85 if risk_forecaster.is_trained else 0.65
    
    return RiskForecastResponse(
        liquidation_risk=risk,
        risk_level=risk_level,
        recommendations=recommendations,
        confidence=confidence
    )


@app.post("/api/v1/apy/trend", response_model=APYTrendResponse)
async def analyze_apy_trend(request: APYTrendRequest):
    """
    Analyze APY trends and predict future rates
    """
    if len(request.historical_apy) < 2:
        raise HTTPException(status_code=400, detail="At least 2 historical APY values required")
    
    # Predict future APY
    predicted_apy = apy_predictor.predict(request.historical_apy)
    
    # Predict trend
    trend = apy_predictor.predict_trend(request.historical_apy)
    
    # Generate recommendation
    current_apy = request.historical_apy[-1] if request.historical_apy else 5.0
    apy_change = predicted_apy - current_apy
    
    if trend == "increasing" and apy_change > 1.0:
        recommendation = f"APY is trending upward. Consider increasing position to capture {apy_change:.2f}% higher returns"
    elif trend == "decreasing" and apy_change < -1.0:
        recommendation = f"APY is trending downward. Consider rebalancing to higher-yield protocols"
    elif trend == "stable":
        recommendation = "APY expected to remain stable. Current position is optimal"
    else:
        recommendation = f"APY trend is {trend}. Monitor for significant changes"
    
    # Confidence based on data quality and model training
    confidence = 0.75 if apy_predictor.is_trained and len(request.historical_apy) >= 7 else 0.60
    
    return APYTrendResponse(
        predicted_apy=round(predicted_apy, 2),
        trend=trend,
        confidence=confidence,
        recommendation=recommendation
    )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)

