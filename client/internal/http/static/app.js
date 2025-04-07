document.addEventListener('DOMContentLoaded', () => {
    // DOM Elements
    const getQuoteBtn = document.getElementById('get-quote');
    const getChallengeBtn = document.getElementById('get-challenge');
    const startLoadBtn = document.getElementById('start-load');
    const stopLoadBtn = document.getElementById('stop-load');
    const quoteText = document.getElementById('quote-text');
    const quoteAuthor = document.getElementById('quote-author');
    const loading = document.getElementById('loading');
    const errorMessage = document.getElementById('error-message');
    const challengeModal = document.getElementById('challenge-modal');
    const closeModal = document.querySelector('.close');
    
    // Stats elements
    const statRequests = document.getElementById('stat-requests');
    const statSuccess = document.getElementById('stat-success');
    const statFailure = document.getElementById('stat-failure');
    const statDifficulty = document.getElementById('stat-difficulty');
    const statDifficultyLevel = document.getElementById('stat-difficulty-level');
    const statScryptR = document.getElementById('stat-scrypt-r');
    const statScryptP = document.getElementById('stat-scrypt-p');
    const statComplexity = document.getElementById('stat-complexity');
    const statAvgTime = document.getElementById('stat-avg-time');
    const statMinTime = document.getElementById('stat-min-time');
    const statMaxTime = document.getElementById('stat-max-time');
    const statLastTime = document.getElementById('stat-last-time');
    const statLoadRequests = document.getElementById('stat-load-requests');
    const statLoadRps = document.getElementById('stat-load-rps');
    
    // Challenge info elements
    const challengeId = document.getElementById('challenge-id');
    const challengeTask = document.getElementById('challenge-task');
    const challengeDifficulty = document.getElementById('challenge-difficulty');
    const challengeN = document.getElementById('challenge-n');
    const challengeR = document.getElementById('challenge-r');
    const challengeP = document.getElementById('challenge-p');
    const challengeKeyLen = document.getElementById('challenge-keylen');
    const challengeComplexity = document.getElementById('challenge-complexity');
    
    // Stats update interval
    let statsInterval = null;
    
    // Initialize
    updateStats();
    
    // Event Listeners
    getQuoteBtn.addEventListener('click', getQuote);
    getChallengeBtn.addEventListener('click', getChallenge);
    startLoadBtn.addEventListener('click', startLoadTest);
    stopLoadBtn.addEventListener('click', stopLoadTest);
    closeModal.addEventListener('click', () => {
        challengeModal.classList.add('hidden');
    });
    
    // Close modal when clicking outside
    window.addEventListener('click', (event) => {
        if (event.target === challengeModal) {
            challengeModal.classList.add('hidden');
        }
    });
    
    // Functions
    async function getQuote() {
        try {
            // Show loading
            quoteText.textContent = '';
            quoteAuthor.textContent = '';
            loading.classList.remove('hidden');
            errorMessage.classList.add('hidden');
            getQuoteBtn.disabled = true;
            
            // Fetch quote
            const response = await fetch('/api/quote');
            const data = await response.json();
            
            // Hide loading
            loading.classList.add('hidden');
            getQuoteBtn.disabled = false;
            
            if (data.success) {
                // Display quote
                quoteText.textContent = `"${data.quote.Text}"`;
                quoteAuthor.textContent = `— ${data.quote.Author}`;
                
                // Update stats
                updateStatsFromData(data.stats);
            } else {
                // Show error
                errorMessage.textContent = data.error;
                errorMessage.classList.remove('hidden');
            }
        } catch (error) {
            // Handle error
            loading.classList.add('hidden');
            getQuoteBtn.disabled = false;
            errorMessage.textContent = `Error: ${error.message}`;
            errorMessage.classList.remove('hidden');
        }
    }
    
    async function getChallenge() {
        try {
            // Show loading
            errorMessage.classList.add('hidden');
            getChallengeBtn.disabled = true;
            
            // Fetch challenge
            const response = await fetch('/api/challenge');
            const data = await response.json();
            
            // Enable button
            getChallengeBtn.disabled = false;
            
            if (data.success && data.challenge) {
                // Display challenge in modal
                challengeId.textContent = data.challenge.challengeId;
                challengeTask.textContent = data.challenge.task;
                challengeDifficulty.textContent = data.challenge.difficultyLevel;
                challengeN.textContent = data.challenge.scryptN;
                challengeR.textContent = data.challenge.scryptR;
                challengeP.textContent = data.challenge.scryptP;
                challengeKeyLen.textContent = data.challenge.keyLen;
                challengeComplexity.textContent = formatComplexity(data.challenge.estimatedComplexity);
                
                // Show modal
                challengeModal.classList.remove('hidden');
                
                // Update stats
                updateStatsFromData(data.stats);
            } else {
                // Show error
                errorMessage.textContent = data.error || 'Failed to get challenge';
                errorMessage.classList.remove('hidden');
            }
        } catch (error) {
            // Handle error
            getChallengeBtn.disabled = false;
            errorMessage.textContent = `Error: ${error.message}`;
            errorMessage.classList.remove('hidden');
        }
    }
    
    async function startLoadTest() {
        try {
            // Disable/enable buttons
            startLoadBtn.disabled = true;
            stopLoadBtn.disabled = false;
            
            // Start load test
            const response = await fetch('/api/load/start', {
                method: 'POST'
            });
            const data = await response.json();
            
            if (!data.success) {
                // Show error
                errorMessage.textContent = data.error || 'Failed to start load test';
                errorMessage.classList.remove('hidden');
                
                // Reset buttons
                startLoadBtn.disabled = false;
                stopLoadBtn.disabled = true;
            } else {
                // Start stats update interval
                if (!statsInterval) {
                    statsInterval = setInterval(updateStats, 1000);
                }
            }
        } catch (error) {
            // Handle error
            errorMessage.textContent = `Error: ${error.message}`;
            errorMessage.classList.remove('hidden');
            
            // Reset buttons
            startLoadBtn.disabled = false;
            stopLoadBtn.disabled = true;
        }
    }
    
    async function stopLoadTest() {
        try {
            // Disable/enable buttons
            stopLoadBtn.disabled = true;
            startLoadBtn.disabled = false;
            
            // Stop load test
            const response = await fetch('/api/load/stop', {
                method: 'POST'
            });
            const data = await response.json();
            
            if (!data.success) {
                // Show error
                errorMessage.textContent = data.error || 'Failed to stop load test';
                errorMessage.classList.remove('hidden');
                
                // Reset buttons
                stopLoadBtn.disabled = false;
                startLoadBtn.disabled = true;
            } else {
                // Clear stats update interval
                if (statsInterval) {
                    clearInterval(statsInterval);
                    statsInterval = null;
                }
                
                // Update stats one last time
                updateStats();
            }
        } catch (error) {
            // Handle error
            errorMessage.textContent = `Error: ${error.message}`;
            errorMessage.classList.remove('hidden');
            
            // Reset buttons
            stopLoadBtn.disabled = false;
            startLoadBtn.disabled = true;
        }
    }
    
    async function updateStats() {
        try {
            const response = await fetch('/api/stats');
            const stats = await response.json();
            
            updateStatsFromData(stats);
        } catch (error) {
            console.error('Error updating stats:', error);
        }
    }
    
    function updateStatsFromData(stats) {
        if (!stats) return;
        
        // Update basic stats display
        statRequests.textContent = stats.requestCount;
        statSuccess.textContent = stats.successCount;
        statFailure.textContent = stats.failureCount;
        
        // Update difficulty stats
        if (stats.lastDifficulty) {
            statDifficulty.textContent = formatNumber(stats.lastDifficulty);
        } else {
            statDifficulty.textContent = 'N/A';
        }
        
        if (stats.lastDifficultyLevel) {
            statDifficultyLevel.textContent = stats.lastDifficultyLevel;
        } else {
            statDifficultyLevel.textContent = 'N/A';
        }
        
        if (stats.lastScryptR) {
            statScryptR.textContent = stats.lastScryptR;
        } else {
            statScryptR.textContent = 'N/A';
        }
        
        if (stats.lastScryptP) {
            statScryptP.textContent = stats.lastScryptP;
        } else {
            statScryptP.textContent = 'N/A';
        }
        
        if (stats.estimatedComplexity) {
            statComplexity.textContent = formatComplexity(stats.estimatedComplexity);
        } else {
            statComplexity.textContent = 'N/A';
        }
        
        // Update performance stats
        if (stats.averageSolveTime) {
            statAvgTime.textContent = `${stats.averageSolveTime.toFixed(2)}s`;
        } else {
            statAvgTime.textContent = 'N/A';
        }
        
        if (stats.minSolveTime) {
            statMinTime.textContent = `${stats.minSolveTime.toFixed(2)}s`;
        } else {
            statMinTime.textContent = 'N/A';
        }
        
        if (stats.maxSolveTime) {
            statMaxTime.textContent = `${stats.maxSolveTime.toFixed(2)}s`;
        } else {
            statMaxTime.textContent = 'N/A';
        }
        
        if (stats.lastSolveTime) {
            statLastTime.textContent = `${stats.lastSolveTime.toFixed(2)}s`;
        } else {
            statLastTime.textContent = 'N/A';
        }
        
        // Update load test stats
        statLoadRequests.textContent = stats.loadTestRequests || 0;
        
        if (stats.loadTestRequestsPerSec) {
            statLoadRps.textContent = stats.loadTestRequestsPerSec.toFixed(2);
        } else {
            statLoadRps.textContent = '0';
        }
        
        // Update button states based on load test status
        if (stats.loadTestActive) {
            startLoadBtn.disabled = true;
            stopLoadBtn.disabled = false;
            
            // Ensure stats interval is running
            if (!statsInterval) {
                statsInterval = setInterval(updateStats, 1000);
            }
        } else {
            startLoadBtn.disabled = false;
            stopLoadBtn.disabled = true;
        }
    }
    
    // Helper function to format large numbers
    function formatNumber(num) {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + 'M';
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + 'K';
        }
        return num.toString();
    }
    
    // Helper function to format complexity values
    function formatComplexity(complexity) {
        if (complexity >= 1e18) {
            return (complexity / 1e18).toFixed(2) + ' Exa';
        } else if (complexity >= 1e15) {
            return (complexity / 1e15).toFixed(2) + ' Peta';
        } else if (complexity >= 1e12) {
            return (complexity / 1e12).toFixed(2) + ' Tera';
        } else if (complexity >= 1e9) {
            return (complexity / 1e9).toFixed(2) + ' Giga';
        } else if (complexity >= 1e6) {
            return (complexity / 1e6).toFixed(2) + ' Mega';
        } else if (complexity >= 1e3) {
            return (complexity / 1e3).toFixed(2) + ' Kilo';
        }
        return complexity.toFixed(2);
    }
});