// Blog page interactive functionality

document.addEventListener('DOMContentLoaded', function() {
    console.log('Blog page loaded successfully!');

    // Demo button functionality
    const demoBtn = document.getElementById('demoBtn');
    const demoOutput = document.getElementById('demoOutput');

    if (demoBtn && demoOutput) {
        demoBtn.addEventListener('click', runDemo);
    }

    // Add smooth scrolling to all links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth'
                });
            }
        });
    });

    // Add animation on scroll
    observeElements();
});

// Demo function
function runDemo() {
    const demoOutput = document.getElementById('demoOutput');
    const demoBtn = document.getElementById('demoBtn');

    // Disable button during demo
    demoBtn.disabled = true;
    demoBtn.textContent = 'Running...';

    // Simulate API call
    setTimeout(() => {
        const demoData = {
            status: 'success',
            message: 'Demo completed successfully!',
            timestamp: new Date().toLocaleString(),
            server: 'Go Web Server',
            protocol: 'HTTP/1.1'
        };

        demoOutput.innerHTML = `
            <h4>Demo Results:</h4>
            <ul>
                <li><strong>Status:</strong> ${demoData.status}</li>
                <li><strong>Message:</strong> ${demoData.message}</li>
                <li><strong>Timestamp:</strong> ${demoData.timestamp}</li>
                <li><strong>Server:</strong> ${demoData.server}</li>
                <li><strong>Protocol:</strong> ${demoData.protocol}</li>
            </ul>
        `;

        demoOutput.classList.add('show');

        // Re-enable button
        demoBtn.disabled = false;
        demoBtn.textContent = 'Run Demo Again';
    }, 1000);
}

// Observe elements for scroll animations
function observeElements() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.style.opacity = '1';
                entry.target.style.transform = 'translateY(0)';
            }
        });
    }, {
        threshold: 0.1
    });

    // Observe blog posts
    document.querySelectorAll('.blog-post').forEach((post, index) => {
        post.style.opacity = '0';
        post.style.transform = 'translateY(20px)';
        post.style.transition = `opacity 0.5s ease ${index * 0.1}s, transform 0.5s ease ${index * 0.1}s`;
        observer.observe(post);
    });
}

// Add dynamic clock to footer
function updateClock() {
    const footer = document.querySelector('.footer p');
    if (footer) {
        const now = new Date();
        const timeString = now.toLocaleTimeString();
        const originalText = footer.textContent.split(' | ')[0];
        footer.textContent = `${originalText} | Current time: ${timeString}`;
    }
}

// Update clock every second
setInterval(updateClock, 1000);
updateClock();

// Log page interactions
console.log('Blog page JavaScript loaded and ready!');
