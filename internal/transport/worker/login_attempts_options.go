package worker

type LoginAttemptOption func(*LoginAttempts)

func WithName(name string) LoginAttemptOption {
    return func(la *LoginAttempts) {
        la.
    }
}
